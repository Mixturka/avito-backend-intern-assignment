package pgxadapter

import "github.com/jackc/pgx/v5"

type row struct {
	pgxRow pgx.Row
}

func NewRow(pgxRow pgx.Row) *row {
	return &row{
		pgxRow: pgxRow,
	}
}

func (r *row) Scan(dest ...any) error {
	return r.pgxRow.Scan(dest...)
}
