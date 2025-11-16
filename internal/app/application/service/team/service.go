package team

import (
	"avito-backend-intern-assignment/internal/app/application/service/user"
	"avito-backend-intern-assignment/internal/app/domain/entity"
	"avito-backend-intern-assignment/pkg/db"
	"context"
	"errors"
	"fmt"
)

type Repository interface {
	db.TransactionalRepository[Repository]
	Create(ctx context.Context, teamName string) error
	Exists(ctx context.Context, teamName string) (bool, error)
}

var (
	ErrTeamUpdate        = errors.New("update team members")
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrTeamNotFound      = errors.New("team not found")
)

type Service struct {
	teamRepo   Repository
	userRepo   user.Repository
	txProvider db.Transactional
}

func NewService(teamRepo Repository, userRepo user.Repository, txProvider db.Transactional) *Service {
	return &Service{
		teamRepo:   teamRepo,
		userRepo:   userRepo,
		txProvider: txProvider,
	}
}

func (s *Service) CreateTeam(ctx context.Context, team entity.Team) error {
	return db.WithTx(ctx, s.txProvider, func(ctx context.Context, tx db.Tx) error {
		txTeamRepo := s.teamRepo.WithDB(tx)
		txUserRepo := s.userRepo.WithDB(tx)

		exists, err := txTeamRepo.Exists(ctx, team.TeamName)
		if err != nil {
			return fmt.Errorf("check team exists: %w", err)
		}

		if !exists {
			if err := txTeamRepo.Create(ctx, team.TeamName); err != nil {
				return fmt.Errorf("create team: %w", err)
			}
		}

		for _, tm := range team.Members {
			userEntity := tm.ToDomainUser(team.TeamName)
			existing, err := txUserRepo.GetByID(ctx, userEntity.UserId)
			if err != nil {
				return fmt.Errorf("get user %s: %w", userEntity.UserId, err)
			}

			if existing == nil {
				if err := txUserRepo.Create(ctx, userEntity); err != nil {
					return fmt.Errorf("create user %s: %w", userEntity.UserId, err)
				}
			} else {
				if err := txUserRepo.Update(ctx, userEntity); err != nil {
					return fmt.Errorf("update user %s: %w", userEntity.UserId, err)
				}
			}
		}

		return nil
	})
}

func (s *Service) GetTeamWithMembers(ctx context.Context, teamName string) (entity.Team, error) {
	exists, err := s.teamRepo.Exists(ctx, teamName)
	if err != nil {
		return entity.Team{}, err
	}
	if !exists {
		return entity.Team{}, ErrTeamNotFound
	}

	users, err := s.userRepo.GetByTeam(ctx, teamName)
	if err != nil {
		return entity.Team{}, err
	}

	members := make([]entity.TeamMember, len(users))
	for i, u := range users {
		members[i] = u.ToDomainTeamMember()
	}

	return entity.Team{
		TeamName: teamName,
		Members:  members,
	}, nil
}
