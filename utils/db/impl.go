package db

import (
	"context"
	"database/sql"
)

type implDb struct {
	db *sql.DB
}

func (i *implDb) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := i.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &implTx{
		tx: tx,
	}, nil
}

func (i *implDb) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := i.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &implRows{
		rows: rows,
	}, nil
}

func (i *implDb) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return &implRow{
		row: i.db.QueryRowContext(ctx, query, args...),
	}
}

type implTx struct {
	tx *sql.Tx
}

func (i *implTx) Commit() error {
	return i.tx.Commit()
}

func (i *implTx) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := i.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return &implRows{
		rows: rows,
	}, nil
}

func (i *implTx) Rollback() error {
	return i.tx.Rollback()
}

type implRows struct {
	rows *sql.Rows
}

func (i *implRows) Close() error {
	return i.rows.Close()
}

func (i *implRows) Next() bool {
	return i.rows.Next()
}

func (i *implRows) Scan(dest ...any) error {
	return i.rows.Scan(dest...)
}

type implRow struct {
	row *sql.Row
}

func (i *implRow) Err() error {
	return i.row.Err()
}

func (i *implRow) Scan(dest ...any) error {
	return i.row.Scan(dest...)
}

func Open(driverName string, dataSourceName string) (DB, error) {
	rawSql, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &implDb{
		db: rawSql,
	}, nil
}
