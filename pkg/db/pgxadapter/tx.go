package pgxadapter

import (
	"avito-backend-intern-assignment/pkg/db"
	"context"

	"github.com/jackc/pgx/v5"
)

type tx struct {
	tx pgx.Tx
}

func (t *tx) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := t.tx.Exec(ctx, sql, args...)
	return err
}

func (t *tx) Query(ctx context.Context, sql string, args ...any) (db.Rows, error) {
	pgxRows, err := t.tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return &rows{
		pgxRows: pgxRows,
	}, nil
}

func (t *tx) QueryRow(ctx context.Context, sql string, args ...any) db.Row {
	return &row{pgxRow: t.tx.QueryRow(ctx, sql, args...)}
}

func (t *tx) Commit(ctx context.Context) error   { return t.tx.Commit(ctx) }
func (t *tx) Rollback(ctx context.Context) error { return t.tx.Rollback(ctx) }
