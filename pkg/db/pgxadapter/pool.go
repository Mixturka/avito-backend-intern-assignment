package pgxadapter

import (
	"avito-backend-intern-assignment/pkg/db"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PoolAdapter struct {
	Pool *pgxpool.Pool
}

func (p *PoolAdapter) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := p.Pool.Exec(ctx, sql, args...)
	return err
}

func (p *PoolAdapter) QueryRow(ctx context.Context, sql string, args ...any) db.Row {
	return &row{pgxRow: p.Pool.QueryRow(ctx, sql, args...)}
}

func (p *PoolAdapter) Query(ctx context.Context, sql string, args ...any) (db.Rows, error) {
	pgxRows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return &rows{
		pgxRows: pgxRows,
	}, nil
}

func (p *PoolAdapter) BeginTx(ctx context.Context) (db.Tx, error) {
	pgxTx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &tx{tx: pgxTx}, nil
}
