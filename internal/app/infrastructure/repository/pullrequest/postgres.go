package pullrequest

import (
	"avito-backend-intern-assignment/internal/app/application/service/pullrequest"
	"avito-backend-intern-assignment/internal/app/domain/entity"
	"avito-backend-intern-assignment/pkg/db"
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type PostgresRepository struct {
	db db.DB
	sb sq.StatementBuilderType
}

func NewPostgresRepository(db db.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PostgresRepository) WithDB(db db.DB) pullrequest.Repository {
	return &PostgresRepository{
		db: db,
		sb: r.sb,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, pr entity.PullRequest) error {
	query, args, err := r.sb.
		Insert("pullrequests").
		Columns("id", "name", "author_id", "status", "created_at", "merged_at").
		Values(pr.PullRequestId, pr.PullRequestName, pr.AuthorId, string(pr.Status), pr.CreatedAt, pr.MergedAt).
		ToSql()
	if err != nil {
		return err
	}

	if err := r.db.Exec(ctx, query, args...); err != nil {
		return err
	}

	for _, reviewer := range pr.AssignedReviewers {
		query, args, err := r.sb.
			Insert("assigned_pr_reviewers").
			Columns("pr_id", "reviewer_id").
			Values(pr.PullRequestId, reviewer).
			ToSql()
		if err != nil {
			return err
		}
		if err := r.db.Exec(ctx, query, args...); err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*entity.PullRequest, error) {
	query, args, err := r.sb.
		Select("id", "name", "author_id", "status", "created_at", "merged_at").
		From("pullrequests").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)

	var pr entity.PullRequest
	var status string
	if err := row.Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &status, &pr.CreatedAt, &pr.MergedAt); err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	pr.Status = entity.PRStatus(status)

	rowsQuery, rowsArgs, err := r.sb.
		Select("reviewer_id").
		From("assigned_pr_reviewers").
		Where(sq.Eq{"pr_id": pr.PullRequestId}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, rowsQuery, rowsArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var rID string
		if err := rows.Scan(&rID); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, rID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	return &pr, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, prID string, status entity.PRStatus, mergedAt *time.Time) error {
	query, args, err := r.sb.
		Update("pullrequests").
		Set("status", string(status)).
		Set("merged_at", mergedAt).
		Where(sq.Eq{"id": prID}).
		ToSql()
	if err != nil {
		return err
	}

	return r.db.Exec(ctx, query, args...)
}

func (r *PostgresRepository) ReplaceReviewer(ctx context.Context, prID string, oldReviewerID string, newReviewerID string) error {
	delQuery, delArgs, _ := r.sb.
		Delete("assigned_pr_reviewers").
		Where(sq.Eq{"pr_id": prID, "reviewer_id": oldReviewerID}).
		ToSql()
	if err := r.db.Exec(ctx, delQuery, delArgs...); err != nil {
		return err
	}

	insQuery, insArgs, _ := r.sb.
		Insert("assigned_pr_reviewers").
		Columns("pr_id", "reviewer_id").
		Values(prID, newReviewerID).
		ToSql()
	return r.db.Exec(ctx, insQuery, insArgs...)
}

func (r *PostgresRepository) GetByAssignedReviewer(ctx context.Context, userID string) ([]entity.PullRequest, error) {
	query, args, err := r.sb.
		Select("pr.id", "pr.name", "pr.author_id", "pr.status", "pr.created_at", "pr.merged_at").
		From("pullrequests pr").
		Join("assigned_pr_reviewers apr ON pr.id = apr.pr_id").
		Where(sq.Eq{"apr.reviewer_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []entity.PullRequest
	for rows.Next() {
		var pr entity.PullRequest
		if err := rows.Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status, &pr.CreatedAt, &pr.MergedAt); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}

func (r *PostgresRepository) BeginTx(ctx context.Context) (db.Tx, error) {
	if transactional, ok := r.db.(db.Transactional); ok {
		return transactional.BeginTx(ctx)
	}
	return nil, fmt.Errorf("underlying database doesn't support transactions")
}
