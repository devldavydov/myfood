package storage

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/devldavydov/myfood/internal/storage/ent"
)

type TxFn func(ctx context.Context, tx *ent.Tx) (any, error)

type DB struct {
	db *ent.Client
}

func (r *DB) Tx(ctx context.Context, fn TxFn) (any, error) {
	// Begin database transaction.
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("db start tx error: %w", err)
	}

	// Rollback on potential panics.
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()

	// Execute user function.
	result, err := fn(ctx, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction.
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("db commit tx error: %w", err)
	}

	return result, nil
}

func (r *DB) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

func NewDB(dbSQL *sql.DB) (*DB, error) {
	drv := entsql.OpenDB(dialect.SQLite, dbSQL)
	db := ent.NewClient(ent.Driver(drv))

	// Run migration
	if err := db.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}
