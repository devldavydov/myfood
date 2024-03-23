package storage

import (
	"context"
	"database/sql"
	"strings"
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

func go_upper(str string) string {
	return strings.ToUpper(str)
}

func NewStorageSQLite(dbFilePath string, logger *zap.Logger) (*StorageSQLite, error) {
	sql.Register("sqlite3_custom", &gsql.SQLiteDriver{
		ConnectHook: func(conn *gsql.SQLiteConn) error {
			if err := conn.RegisterFunc("go_upper", go_upper, false); err != nil {
				return err
			}
			return nil
		},
	})

	db, err := sql.Open("sqlite3_custom", dbFilePath)
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

func (r *StorageSQLite) GetFood(ctx context.Context, key string) (*Food, error) {
	var f Food
	err := r.db.
		QueryRowContext(ctx, _sqlGetFood, key).
		Scan(&f.Key, &f.Name, &f.Brand, &f.Cal100, &f.Prot100, &f.Fat100, &f.Carb100, &f.Comment)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrFoodNotFound
	case err != nil:
		return nil, err
	}

	return &f, nil
}

func (r *StorageSQLite) SetFood(ctx context.Context, food *Food) error {
	if !food.Validate() {
		return ErrFoodInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetFood,
		food.Key,
		food.Name,
		food.Brand,
		food.Cal100,
		food.Prot100,
		food.Fat100,
		food.Carb100,
		food.Comment,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *StorageSQLite) GetFoodList(ctx context.Context) ([]Food, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetFoodList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Food
	for rows.Next() {
		var f Food
		err = rows.Scan(&f.Key, &f.Name, &f.Brand, &f.Cal100, &f.Prot100, &f.Fat100, &f.Carb100, &f.Comment)
		if err != nil {
			return nil, err
		}

		list = append(list, f)
	}

	if len(list) == 0 {
		return nil, ErrFoodEmptyList
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *StorageSQLite) FindFood(ctx context.Context, pattern string) ([]Food, error) {
	rows, err := r.db.QueryContext(ctx, _sqFindFood, strings.ToUpper(pattern))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Food
	for rows.Next() {
		var f Food
		err = rows.Scan(&f.Key, &f.Name, &f.Brand, &f.Cal100, &f.Prot100, &f.Fat100, &f.Carb100, &f.Comment)
		if err != nil {
			return nil, err
		}

		list = append(list, f)
	}

	if len(list) == 0 {
		return nil, ErrFoodEmptyList
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *StorageSQLite) DeleteFood(ctx context.Context, key string) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteFood, key)
	return err
}

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
