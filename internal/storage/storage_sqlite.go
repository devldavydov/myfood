package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/devldavydov/myfood/internal/storage/ent"
	"github.com/devldavydov/myfood/internal/storage/ent/usersettings"
	"github.com/devldavydov/myfood/internal/storage/ent/weight"
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
	dbEnt  *DB
	logger *zap.Logger
}

var _ Storage = (*StorageSQLite)(nil)

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

	//
	// Open DB entgo.
	//

	dbEnt, err := NewDB(dbFilePath)
	if err != nil {
		return nil, err
	}

	stg := &StorageSQLite{db: db, dbEnt: dbEnt, logger: logger}
	if err := stg.init(); err != nil {
		return nil, err
	}

	// Run migrations from old to new (remove after all done)
	// if err := stg.migrateWeight(); err != nil {
	// 	return nil, err
	// }

	return stg, nil
}

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

func (r *StorageSQLite) SetFoodComment(ctx context.Context, key, comment string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get rowid to check existence
	var rowid int64
	err = tx.QueryRowContext(ctx, _sqlFoodRowid, key).Scan(&rowid)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrFoodNotFound
	case err != nil:
		return err
	}

	// Update
	_, err = tx.ExecContext(ctx, _sqlSetFoodComment, comment, key)
	if err != nil {
		return err
	}

	return tx.Commit()
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
	if err != nil {
		var errSql gsql.Error
		if errors.As(err, &errSql) && errSql.Error() == _errForeignKey {
			return ErrFoodIsUsed
		}
		return err
	}
	return nil
}

//
// Weight
//

func (r *StorageSQLite) GetWeightList(ctx context.Context, userID int64, from, to time.Time) ([]Weight, error) {
	res, err := r.dbEnt.Tx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Weight.
			Query().
			Where(
				weight.Userid(userID),
				weight.TimestampGTE(from),
				weight.TimestampLTE(to),
			).
			Order(
				weight.ByTimestamp(),
			).
			All(ctx)
	})
	if err != nil {
		return nil, err
	}

	weLst, _ := res.([]*ent.Weight)
	if len(weLst) == 0 {
		return nil, ErrWeightEmptyList
	}

	wLst := make([]Weight, 0, len(weLst))
	for _, w := range weLst {
		wLst = append(wLst, Weight{Timestamp: w.Timestamp, Value: w.Value})
	}

	return wLst, nil
}

func (r *StorageSQLite) SetWeight(ctx context.Context, userID int64, weight *Weight) error {
	if !weight.Validate() {
		return ErrWeightInvalid
	}

	_, err := r.dbEnt.Tx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Weight.
			Create().
			SetUserid(userID).
			SetTimestamp(weight.Timestamp).
			SetValue(weight.Value).
			OnConflict().
			UpdateNewValues().
			ID(ctx)
	})

	return err
}

func (r *StorageSQLite) DeleteWeight(ctx context.Context, userID int64, timestamp time.Time) error {
	_, err := r.dbEnt.Tx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Weight.
			Delete().
			Where(
				weight.Userid(userID),
				weight.Timestamp(timestamp),
			).
			Exec(ctx)
	})
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

func (r *StorageSQLite) GetJournalReport(ctx context.Context, userID int64, from, to int64) ([]JournalReport, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetJournalReport, userID, from, to)
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
			&jd.FoodKey,
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

func (r *StorageSQLite) GetJournalStats(ctx context.Context, userID int64, from, to int64) ([]JournalStats, error) {
	rows, err := r.db.QueryContext(ctx, _sqlGetJournalStats, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []JournalStats
	for rows.Next() {
		var js JournalStats
		err = rows.Scan(
			&js.Timestamp,
			&js.TotalCal,
			&js.TotalProt,
			&js.TotalFat,
			&js.TotalCarb,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, js)
	}

	if len(list) == 0 {
		return nil, ErrJournalStatsEmpty
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
	res, err := r.dbEnt.Tx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.UserSettings.
			Query().
			Where(usersettings.Userid(userID)).
			First(ctx)
	})

	if err != nil {
		var notFound *ent.NotFoundError
		if errors.As(err, &notFound) {
			return nil, ErrUserSettingsNotFound
		}

		return nil, err
	}

	us, _ := res.(*ent.UserSettings)

	return &UserSettings{CalLimit: us.CalLimit}, nil
}

func (r *StorageSQLite) SetUserSettings(ctx context.Context, userID int64, settings *UserSettings) error {
	if !settings.Validate() {
		return ErrUserSettingsInvalid
	}

	_, err := r.dbEnt.Tx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.UserSettings.
			Create().
			SetUserid(userID).
			SetCalLimit(settings.CalLimit).
			OnConflict().
			UpdateNewValues().
			ID(ctx)
	})

	return err
}

//
//
//

func (r *StorageSQLite) Close() error {
	if r.dbEnt == nil {
		return nil
	}

	return r.dbEnt.Close()
}

func (r *StorageSQLite) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), _databaseInitTimeout)
	defer cancel()

	for _, createTbl := range []string{
		_sqlCreateTableFood,
		_sqlCreateTableJournal} {
		_, err := r.db.ExecContext(ctx, createTbl)
		if err != nil {
			return err
		}

	}

	return nil
}

func (r *StorageSQLite) migrateWeight() error {
	// Get from old table
	ctx := context.Background()
	rows, err := r.db.QueryContext(ctx, `
	SELECT userid, timestamp, value
	FROM weight2
	`)

	if err != nil {
		return err
	}
	defer rows.Close()

	type row struct {
		userid    int64
		timestamp int64
		value     float64
	}

	var list []row
	for rows.Next() {
		var f row
		err = rows.Scan(&f.userid, &f.timestamp, &f.value)
		if err != nil {
			return err
		}

		list = append(list, f)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	// Save in new format
	_, err = r.dbEnt.Tx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		for _, i := range list {
			_, err := tx.Weight.
				Create().
				SetUserid(i.userid).
				SetValue(i.value).
				SetTimestamp(time.Unix(i.timestamp, 0)).
				Save(ctx)

			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	// Delete from table
	_, err = r.db.ExecContext(ctx, `
	DELETE FROM weight2;
	`)

	return err
}
