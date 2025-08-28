package repositories

import (
	"context"

	"github.com/jackc/pgx"
)

type PGXTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgx.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
