package storage

import (
	"context"
	"fmt"

	"github.com/devldavydov/myfood/internal/storage/ent"
)

type TxFn func(ctx context.Context, tx *ent.Tx) (any, error)

type DB struct {
	db   *ent.Client
	path string
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

func NewDB(dbFilePath string) (*DB, error) {
	// Format url
	url := fmt.Sprintf(
		"file:%s?mode=rwc&_timeout=5000&_fk=1&_sync=1&_journal=wal",
		dbFilePath,
	)

	// Open DB
	db, err := ent.Open("sqlite3", url)
	if err != nil {
		return nil, fmt.Errorf("open database with %s: %w", url, err)
	}

	// Run migration
	if err := db.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("database %s migration error: %w", url, err)
	}

	return &DB{
		db:   db,
		path: dbFilePath,
	}, nil
}
