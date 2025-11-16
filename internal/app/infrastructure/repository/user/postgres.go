package user

import (
	"avito-backend-intern-assignment/internal/app/application/service/user"
	"avito-backend-intern-assignment/internal/app/domain/entity"
	"avito-backend-intern-assignment/pkg/db"
	"context"
	"fmt"

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

func (r *PostgresRepository) WithDB(db db.DB) user.Repository {
	return &PostgresRepository{
		db: db,
		sb: r.sb,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, user entity.User) error {
	query, args, err := r.sb.
		Insert("users").
		Columns("id", "username", "is_active", "team_name").
		Values(user.UserId, user.Username, user.IsActive, user.TeamName).
		ToSql()
	if err != nil {
		return err
	}
	err = r.db.Exec(ctx, query, args...)

	return err
}

func (r *PostgresRepository) Update(ctx context.Context, user entity.User) error {
	query, args, err := r.sb.
		Update("users").
		Set("username", user.Username).
		Set("is_active", user.IsActive).
		Set("team_name", user.TeamName).
		Where(sq.Eq{"id": user.UserId}).
		ToSql()
	if err != nil {
		return err
	}
	err = r.db.Exec(ctx, query, args...)

	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	query, args, err := r.sb.
		Select("id", "username", "is_active", "team_name").
		From("users").
		Where(sq.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	var u entity.User
	if err := row.Scan(&u.UserId, &u.Username, &u.IsActive, &u.TeamName); err != nil {
		return nil, nil
	}

	return &u, nil
}

func (r *PostgresRepository) GetByTeam(ctx context.Context, teamName string) ([]entity.User, error) {
	query, args, err := r.sb.
		Select("id", "username", "is_active", "team_name").
		From("users").
		Where(sq.Eq{"team_name": teamName}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.UserId, &u.Username, &u.IsActive, &u.TeamName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresRepository) BeginTx(ctx context.Context) (db.Tx, error) {
	if transactional, ok := r.db.(db.Transactional); ok {
		return transactional.BeginTx(ctx)
	}
	return nil, fmt.Errorf("underlying database doesn't support transactions")
}
