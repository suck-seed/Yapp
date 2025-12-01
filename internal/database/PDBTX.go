package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// A common interface for both ConnWrapper + TxWrapper
type DBRunner interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Transaction Interface
type PGXTX interface {
	DBRunner
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Implements PGXTX
type TxWrapper struct {
	tx pgx.Tx
}

func NewTxWrapper(tx pgx.Tx) *TxWrapper {
	return &TxWrapper{tx: tx}
}

func (w *TxWrapper) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return w.tx.Exec(ctx, sql, args...)
}

func (w *TxWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return w.tx.Query(ctx, sql, args...)
}

func (w *TxWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return w.tx.QueryRow(ctx, sql, args...)
}

func (w *TxWrapper) Commit(ctx context.Context) error {
	return w.tx.Commit(ctx)
}

func (w *TxWrapper) Rollback(ctx context.Context) error {
	return w.tx.Rollback(ctx)
}

// Implements DBRunner
type ConnWrapper struct {
	conn *pgxpool.Conn
}

func NewConnWrapper(conn *pgxpool.Conn) *ConnWrapper {
	return &ConnWrapper{conn}
}

func (w *ConnWrapper) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return w.conn.Exec(ctx, sql, args...)
}

func (w *ConnWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return w.conn.Query(ctx, sql, args...)
}

func (w *ConnWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return w.conn.QueryRow(ctx, sql, args...)
}
