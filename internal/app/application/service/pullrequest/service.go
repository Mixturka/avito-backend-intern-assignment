package pullrequest

import (
	"avito-backend-intern-assignment/internal/app/application/service/user"
	"avito-backend-intern-assignment/internal/app/domain/entity"
	"avito-backend-intern-assignment/pkg/db"
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"time"
)

type Repository interface {
	db.TransactionalRepository[Repository]
	Create(ctx context.Context, pr entity.PullRequest) error
	GetByID(ctx context.Context, id string) (*entity.PullRequest, error)
	UpdateStatus(ctx context.Context, prID string, status entity.PRStatus, mergedAt *time.Time) error
	ReplaceReviewer(ctx context.Context, prID string, oldReviewerID string, newReviewerID string) error
	GetByAssignedReviewer(ctx context.Context, userID string) ([]entity.PullRequest, error)
}

var (
	ErrPullRequestExists   = errors.New("pull request already exists")
	ErrPullRequestNotFound = errors.New("pull request not found")
	ErrAuthorNotFound      = errors.New("author not found")
	ErrTeamNotFound        = errors.New("team not found")
	ErrReassignViolation   = errors.New("reassigning reviewer violates domain rules")
	ErrNotEnoughReviewers  = errors.New("not enough active reviewers in the team")
)

type Service struct {
	prRepo     Repository
	userRepo   user.Repository
	txProvider db.Transactional
}

func NewService(prRepo Repository, userRepo user.Repository, txProvider db.Transactional) *Service {
	return &Service{
		prRepo:     prRepo,
		userRepo:   userRepo,
		txProvider: txProvider,
	}
}

func (s *Service) getRandomTeamReviewers(ctx context.Context, userRepo user.Repository, teamName string, n int, excludedUserIDs ...string) ([]string, error) {
	teamMembers, err := userRepo.GetByTeam(ctx, teamName)
	if err != nil {
		log.Printf("ERROR: Failed to get team members for team '%s': %v", teamName, err)
		return nil, fmt.Errorf("get team members: %w", err)
	}

	excludedSet := make(map[string]bool)
	for _, id := range excludedUserIDs {
		excludedSet[id] = true
	}

	potentialReviewers := make([]string, 0)
	for _, member := range teamMembers {
		if member.IsActive && !excludedSet[member.UserId] {
			potentialReviewers = append(potentialReviewers, member.UserId)
		}
	}

	if len(potentialReviewers) == 0 {
		log.Printf("ERROR: No potential reviewers found for team '%s'. Team members: %d, Excluded: %v",
			teamName, len(teamMembers), excludedUserIDs)
		return nil, ErrNotEnoughReviewers
	}

	selected := s.selectRandomReviewers(potentialReviewers, n)

	return selected, nil
}

func (s *Service) selectRandomReviewers(reviewers []string, n int) []string {
	if len(reviewers) <= n {
		shuffled := make([]string, len(reviewers))
		copy(shuffled, reviewers)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})

		return shuffled
	}

	selected := make([]string, 0, n)
	indexes := rand.Perm(len(reviewers))
	for i := range n {
		selected = append(selected, reviewers[indexes[i]])
	}

	return selected
}

func (s *Service) Create(ctx context.Context, pr entity.PullRequest) (*entity.PullRequest, error) {
	existing, err := s.prRepo.GetByID(ctx, pr.PullRequestId)
	if err != nil {
		log.Printf("ERROR: Failed to check if PR exists (ID: %s): %v", pr.PullRequestId, err)
		return nil, fmt.Errorf("check pr exists: %w", err)
	}
	if existing != nil {
		log.Printf("ERROR: PR already exists (ID: %s)", pr.PullRequestId)
		return nil, ErrPullRequestExists
	}

	author, err := s.userRepo.GetByID(ctx, pr.AuthorId)
	if err != nil {
		log.Printf("ERROR: Failed to get author (ID: %s): %v", pr.AuthorId, err)
		return nil, fmt.Errorf("get author: %w", err)
	}
	if author == nil {
		log.Printf("ERROR: Author not found (ID: %s)", pr.AuthorId)
		return nil, ErrAuthorNotFound
	}

	selectedReviewers, err := s.getRandomTeamReviewers(ctx, s.userRepo, author.TeamName, 2, pr.AuthorId)
	if err != nil {
		log.Printf("ERROR: Failed to select reviewers for PR %s: %v", pr.PullRequestId, err)
		return nil, fmt.Errorf("select reviewers: %w", err)
	}

	pr.AssignedReviewers = selectedReviewers
	pr.Status = entity.PullRequestStatusOPEN
	createdAt := time.Now().UTC()
	pr.CreatedAt = &createdAt

	if err := s.prRepo.Create(ctx, pr); err != nil {
		log.Printf("ERROR: Failed to create PR in database (ID: %s): %v", pr.PullRequestId, err)
		return nil, fmt.Errorf("create pr in database: %w", err)
	}

	return &pr, nil
}

func (s *Service) MarkMerged(ctx context.Context, prID string) (*entity.PullRequest, error) {
	log.Printf("Marking PR as merged: ID=%s", prID)

	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		log.Printf("ERROR: Failed to get PR for merging (ID: %s): %v", prID, err)
		return nil, fmt.Errorf("get pr: %w", err)
	}
	if pr == nil {
		log.Printf("ERROR: PR not found for merging (ID: %s)", prID)
		return nil, ErrPullRequestNotFound
	}

	if pr.Status == entity.PullRequestStatusMERGED {
		log.Printf("PR %s is already merged", prID)
		return pr, nil
	}

	now := time.Now().UTC()
	log.Printf("Updating PR %s status to MERGED at %v", prID, now)

	if err := s.prRepo.UpdateStatus(ctx, prID, entity.PullRequestStatusMERGED, &now); err != nil {
		log.Printf("ERROR: Failed to update PR status (ID: %s): %v", prID, err)
		return nil, fmt.Errorf("update pr status: %w", err)
	}

	pr.Status = entity.PullRequestStatusMERGED
	pr.MergedAt = &now

	log.Printf("Successfully marked PR %s as merged", prID)
	return pr, nil
}

func (s *Service) ReassignReviewer(ctx context.Context, prID string, oldReviewerID string) (*entity.PullRequest, string, error) {
	log.Printf("Reassigning reviewer: PR=%s, OldReviewer=%s", prID, oldReviewerID)

	var updatedPR *entity.PullRequest
	var newReviewerID string

	err := db.WithTx(ctx, s.txProvider, func(ctx context.Context, tx db.Tx) error {
		txRepo := s.prRepo.WithDB(tx)
		txUserRepo := s.userRepo.WithDB(tx)

		pr, err := txRepo.GetByID(ctx, prID)
		if err != nil {
			log.Printf("ERROR: Failed to get PR for reassignment (ID: %s): %v", prID, err)
			return fmt.Errorf("get pr: %w", err)
		}
		if pr == nil {
			log.Printf("ERROR: PR not found for reassignment (ID: %s)", prID)
			return ErrPullRequestNotFound
		}
		log.Printf("PR found: ID=%s, Author=%s, Reviewers=%v", pr.PullRequestId, pr.AuthorId, pr.AssignedReviewers)

		found := slices.Contains(pr.AssignedReviewers, oldReviewerID)
		if !found {
			log.Printf("ERROR: Old reviewer %s is not assigned to PR %s. Current reviewers: %v",
				oldReviewerID, prID, pr.AssignedReviewers)
			return ErrReassignViolation
		}

		oldUser, err := txUserRepo.GetByID(ctx, oldReviewerID)
		if err != nil {
			log.Printf("ERROR: Failed to get old reviewer (ID: %s): %v", oldReviewerID, err)
			return fmt.Errorf("get old reviewer: %w", err)
		}
		if oldUser == nil {
			log.Printf("ERROR: Old reviewer not found (ID: %s)", oldReviewerID)
			return user.ErrUserNotFound
		}

		excludedUsers := make([]string, 0, len(pr.AssignedReviewers)+1)
		excludedUsers = append(excludedUsers, pr.AssignedReviewers...)
		excludedUsers = append(excludedUsers, pr.AuthorId)

		log.Printf("Excluding users for replacement: %v", excludedUsers)

		candidateReviewers, err := s.getRandomTeamReviewers(ctx, txUserRepo, oldUser.TeamName, 1, excludedUsers...)
		if err != nil {
			log.Printf("ERROR: Failed to get replacement reviewer for PR %s: %v", prID, err)
			return fmt.Errorf("get replacement reviewer: %w", err)
		}

		if len(candidateReviewers) == 0 {
			log.Printf("ERROR: No candidate reviewers found for reassignment in team %s", oldUser.TeamName)
			return ErrReassignViolation
		}

		newReviewerID = candidateReviewers[0]
		log.Printf("Selected new reviewer: %s", newReviewerID)

		if err := txRepo.ReplaceReviewer(ctx, prID, oldReviewerID, newReviewerID); err != nil {
			log.Printf("ERROR: Failed to replace reviewer in database (PR: %s, Old: %s, New: %s): %v",
				prID, oldReviewerID, newReviewerID, err)
			return fmt.Errorf("replace reviewer: %w", err)
		}

		for i, r := range pr.AssignedReviewers {
			if r == oldReviewerID {
				pr.AssignedReviewers[i] = newReviewerID
				break
			}
		}

		updatedPR = pr
		log.Printf("Successfully reassigned reviewer: PR=%s, Old=%s, New=%s", prID, oldReviewerID, newReviewerID)
		return nil
	})
	if err != nil {
		log.Printf("ERROR: Transaction failed for reassignment (PR: %s): %v", prID, err)
	}

	return updatedPR, newReviewerID, err
}

func (s *Service) GetPRsByReviewer(ctx context.Context, userID string) (string, []entity.PullRequest, error) {
	log.Printf("Getting PRs for reviewer: %s", userID)

	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get user for PR lookup (ID: %s): %v", userID, err)
		return "", nil, user.ErrUserNotFound
	}
	if u == nil {
		log.Printf("ERROR: User not found for PR lookup (ID: %s)", userID)
		return "", nil, user.ErrUserNotFound
	}
	log.Printf("User found: ID=%s, Team=%s", u.UserId, u.TeamName)

	prs, err := s.prRepo.GetByAssignedReviewer(ctx, userID)
	if err != nil {
		log.Printf("ERROR: Failed to get PRs for reviewer %s: %v", userID, err)
		return "", nil, fmt.Errorf("get PRs for reviewer: %w", err)
	}

	log.Printf("Found %d PRs for reviewer %s", len(prs), userID)
	return userID, prs, nil
}
