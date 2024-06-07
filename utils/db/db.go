package db

import (
	"context"
	"database/sql"
)

type (
	DB interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
		QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...any) Row
	}

	Tx interface {
		Commit() error
		QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
		Rollback() error
	}

	Rows interface {
		Close() error
		Next() bool
		Scan(dest ...any) error
	}

	Row interface {
		Err() error
		Scan(dest ...any) error
	}
)
