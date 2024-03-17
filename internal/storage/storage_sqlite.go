package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	gsql "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const (
	_databaseInitTimeout = 10 * time.Second

	_constraintUniqueWeight = "UNIQUE constraint failed: weight.userid, weight.timestamp"
)

type StorageSQLite struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewStorageSQLite(dbFilePath string, logger *zap.Logger) (*StorageSQLite, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	stg := &StorageSQLite{db: db, logger: logger}
	if err := stg.init(); err != nil {
		return nil, err
	}

	return stg, nil
}

//
// Food
//

//
// Weight
//

func (r *StorageSQLite) GetWeight(ctx context.Context, userID int64, timestamp int64) (*Weight, error) {
	var w Weight
	err := r.db.QueryRowContext(ctx, _sqlFindWeight, userID, timestamp).Scan(&w.Timestamp, &w.Value)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrWeightNotFound
	case err != nil:
		return nil, err
	}

	return &w, nil
}

func (r *StorageSQLite) CreateWeight(ctx context.Context, userID int64, weight *Weight) error {
	if !weight.Validate() {
		return ErrInvalidWeight
	}

	_, err := r.db.ExecContext(ctx, _sqlCreateWeight, userID, weight.Timestamp, weight.Value)
	if err != nil {
		var dbErr gsql.Error
		if !errors.As(err, &dbErr) {
			return err
		}

		if dbErr.Error() == _constraintUniqueWeight {
			return ErrWeightAlreadyExists
		}

		return err
	}

	return nil
}

//
//
//

func (r *StorageSQLite) Close() {
	if r.db == nil {
		return
	}

	if err := r.db.Close(); err != nil {
		r.logger.Error("db close error", zap.Error(err))
	}
}

func (r *StorageSQLite) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseInitTimeout)
	defer cancel()

	for _, createTbl := range []string{_sqlCreateTableFood, _sqlCreateTableWeight} {
		_, err := r.db.ExecContext(ctx, createTbl)
		if err != nil {
			return err
		}

	}

	return nil
}
