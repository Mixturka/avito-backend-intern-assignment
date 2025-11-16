package pgxadapter

import "github.com/jackc/pgx/v5"

type rows struct {
	pgxRows pgx.Rows
}

func NewRows(pgxRows pgx.Rows) *rows {
	return &rows{
		pgxRows: pgxRows,
	}
}

func (r *rows) Next() bool             { return r.pgxRows.Next() }
func (r *rows) Scan(dest ...any) error { return r.pgxRows.Scan(dest...) }
func (r *rows) Close()                 { r.pgxRows.Close() }
func (r *rows) Err() error             { return r.pgxRows.Err() }
