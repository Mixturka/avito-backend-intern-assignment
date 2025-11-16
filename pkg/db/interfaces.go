package db

import (
	"context"
)

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

type QueryExecutor interface {
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
}

type CommandExecutor interface {
	Exec(ctx context.Context, sql string, args ...any) error
}

type Transactional interface {
	BeginTx(ctx context.Context) (Tx, error)
}

type DB interface {
	QueryExecutor
	CommandExecutor
}

type TransactionalDB interface {
	DB
	Transactional
}

type Tx interface {
	DB
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TransactionalRepository[T any] interface {
	// WithDB предполагает создание копии репозитория с заменой внутреннего DB на другой DB.
	// Предлагается использовать метод для подмены на tx при совместном использовании
	// с db.WithTx(...), если необходима транзакционность
	WithDB(db DB) T
	Transactional
}

func WithTx(ctx context.Context, repo Transactional, fn func(ctx context.Context, tx Tx) error) error {
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
