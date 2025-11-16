package team

import (
	"avito-backend-intern-assignment/internal/app/application/service/team"
	"avito-backend-intern-assignment/pkg/db"
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

func (r *PostgresRepository) WithDB(db db.DB) team.Repository {
	return &PostgresRepository{
		db: db,
		sb: r.sb,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, teamName string) error {
	query, args, err := r.sb.
		Insert("teams").
		Columns("team_name").
		Values(teamName).
		ToSql()
	if err != nil {
		return err
	}
	err = r.db.Exec(ctx, query, args...)

	return err
}

func (r *PostgresRepository) Exists(ctx context.Context, teamName string) (bool, error) {
	query, args, err := r.sb.
		Select("1").
		From("teams").
		Where(sq.Eq{"team_name": teamName}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	row := r.db.QueryRow(ctx, query, args...)

	var dummy int
	err = row.Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *PostgresRepository) BeginTx(ctx context.Context) (db.Tx, error) {
	if transactional, ok := r.db.(db.Transactional); ok {
		return transactional.BeginTx(ctx)
	}
	return nil, fmt.Errorf("underlying database doesn't support transactions")
}
