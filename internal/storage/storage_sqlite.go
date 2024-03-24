package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	gsql "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const (
	_databaseInitTimeout = 10 * time.Second

	_customDriverName = "sqlite3_custom"
	_errForeignKey    = "FOREIGN KEY constraint failed"
)

type StorageSQLite struct {
	db     *sql.DB
	logger *zap.Logger
}

func go_upper(str string) string {
	return strings.ToUpper(str)
}

func NewStorageSQLite(dbFilePath string, logger *zap.Logger) (*StorageSQLite, error) {
	//
	// Driver register (check registration twice).
	//

	if !isDriverRegistered(_customDriverName) {
		sql.Register(_customDriverName, &gsql.SQLiteDriver{
			ConnectHook: func(conn *gsql.SQLiteConn) error {
				if err := conn.RegisterFunc("go_upper", go_upper, false); err != nil {
					return err
				}
				return nil
			},
		})
	}

	//
	// Open DB.
	//

	db, err := sql.Open(_customDriverName, dbFilePath+"?_foreign_keys=1")
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
// Journal
//

func (r *StorageSQLite) SetJournal(ctx context.Context, userID int64, journal *Journal) error {
	if !journal.Validate() {
		return ErrJournalInvalid
	}

	_, err := r.db.ExecContext(ctx, _sqlSetJournal, userID, journal.Timestamp, journal.Meal, journal.FoodKey, journal.FoodWeight)
	if err != nil {
		var errSql gsql.Error
		if errors.As(err, &errSql) && errSql.Error() == _errForeignKey {
			return ErrJournalInvalidFood
		}
		return err
	}

	return nil
}

func (r *StorageSQLite) DeleteJournal(ctx context.Context, userID int64, timestamp int64, meal Meal, foodkey string) error {
	_, err := r.db.ExecContext(ctx, _sqlDeleteJournal, userID, timestamp, meal, foodkey)
	return err
}

func (r *StorageSQLite) GetJournalForPeriod(ctx context.Context, userID int64, from, to int64) ([]JournalReport, error) {
	return r.processJournalReport(ctx, _sqlGetJournalForPeriod, userID, from, to)
}

func (r *StorageSQLite) GetJournalForPeriodAndMeal(ctx context.Context, userID int64, from, to int64, meal Meal) ([]JournalReport, error) {
	return r.processJournalReport(ctx, _sqlGetJournalForPeriodAndMeal, userID, from, to, meal)
}

func (r *StorageSQLite) processJournalReport(ctx context.Context, query string, args ...any) ([]JournalReport, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []JournalReport
	for rows.Next() {
		var jd JournalReport
		err = rows.Scan(
			&jd.Timestamp,
			&jd.Meal,
			&jd.FoodName,
			&jd.FoodBrand,
			&jd.FoodWeight,
			&jd.Cal,
			&jd.Prot,
			&jd.Fat,
			&jd.Carb,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, jd)
	}

	if len(list) == 0 {
		return nil, ErrJournalReportEmpty
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
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

	for _, createTbl := range []string{
		_sqlCreateTableFood,
		_sqlCreateTableJournal,
		_sqlCreateTableWeight,
		_sqlCreateTableUserSettings} {
		_, err := r.db.ExecContext(ctx, createTbl)
		if err != nil {
			return err
		}

	}

	return nil
}
