package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

var _ Storage = (*StorageSQLite)(nil)

//
// Food
//

//
// Weight
//

func (r *StorageSQLite) GetWeightList(ctx context.Context, userID int64, from, to int64) ([]Weight, error) {
	rows, err := r.db.QueryContext(ctx, _sqlWeightList, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Weight
	for rows.Next() {
		var w Weight
		err = rows.Scan(&w.Timestamp, &w.Value)
		if err != nil {
			return nil, err
		}

		list = append(list, w)
	}

	if len(list) == 0 {
		return nil, ErrWeightEmptyList
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *StorageSQLite) SetWeight(ctx context.Context, userID int64, weight *Weight) error {
	if !weight.Validate() {
		return ErrWeightInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetWeight, userID, weight.Timestamp, weight.Value)
	if err != nil {
		return err
	}

	return nil
}

func (r *StorageSQLite) DeleteWeight(ctx context.Context, userID, timestamp int64) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteWeight, userID, timestamp)
	return err
}

//
// UserSettings
//

func (r *StorageSQLite) GetUserSettings(ctx context.Context, userID int64) (*UserSettings, error) {
	var us UserSettings
	err := r.db.QueryRowContext(ctx, _sqlGetUserSettings, userID).Scan(&us.CalLimit)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrUserSettingsNotFound
	case err != nil:
		return nil, err
	}

	return &us, nil
}

func (r *StorageSQLite) SetUserSettings(ctx context.Context, userID int64, settings *UserSettings) error {
	if !settings.Validate() {
		return ErrUserSettingsInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetUserSettings, userID, settings.CalLimit)
	if err != nil {
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

	for _, createTbl := range []string{_sqlCreateTableFood, _sqlCreateTableWeight, _sqlCreateTableUserSettings} {
		_, err := r.db.ExecContext(ctx, createTbl)
		if err != nil {
			return err
		}

	}

	return nil
}
