package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/devldavydov/myfood/internal/storage/ent"
	"github.com/devldavydov/myfood/internal/storage/ent/food"
	"github.com/devldavydov/myfood/internal/storage/ent/usersettings"
	"github.com/devldavydov/myfood/internal/storage/ent/weight"
	gsql "github.com/mattn/go-sqlite3"
)

const (
	_databaseInitTimeout = 10 * time.Second

	_customDriverName = "sqlite3_custom"
	_errForeignKey    = "FOREIGN KEY constraint failed"
)

type TxFn func(ctx context.Context, tx *ent.Tx) (any, error)

type StorageSQLite struct {
	dbSQL *sql.DB // Remove
	db    *ent.Client
}

var _ Storage = (*StorageSQLite)(nil)

func go_upper(str string) string {
	return strings.ToUpper(str)
}

func NewStorageSQLite(dbFilePath string) (*StorageSQLite, error) {
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

	// Format url
	url := fmt.Sprintf(
		"file:%s?mode=rwc&_timeout=5000&_fk=1&_sync=1&_journal=wal",
		dbFilePath,
	)

	//
	// Open DB.
	//

	dbSQL, err := sql.Open(_customDriverName, url)
	if err != nil {
		return nil, err
	}

	//
	// Open DB entgo.
	//

	drv := entsql.OpenDB(dialect.SQLite, dbSQL)
	dbEnt := ent.NewClient(ent.Driver(drv))

	// Run migration
	if err := dbEnt.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	stg := &StorageSQLite{db: dbEnt}

	// Run migrations from old to new (remove after all done)
	// if err := stg.migrateWeight(); err != nil {
	// 	return nil, err
	// }

	return stg, nil
}

func (r *StorageSQLite) doTx(ctx context.Context, fn TxFn) (any, error) {
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

func (r *StorageSQLite) Close() error {
	if r.db == nil {
		return nil
	}

	return r.db.Close()
}

//
// Food
//

func (r *StorageSQLite) GetFood(ctx context.Context, key string) (*Food, error) {
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Food.
			Query().
			Where(food.Key(key)).
			First(ctx)
	})

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrFoodNotFound
		}
		return nil, err
	}

	ef, _ := res.(*ent.Food)
	return foodFromEntFood(ef), nil
}

func (r *StorageSQLite) SetFood(ctx context.Context, food *Food) error {
	if !food.Validate() {
		return ErrFoodInvalid
	}

	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Food.
			Create().
			SetKey(food.Key).
			SetName(food.Name).
			SetBrand(food.Brand).
			SetCal100(food.Cal100).
			SetProt100(food.Prot100).
			SetFat100(food.Fat100).
			SetCarb100(food.Carb100).
			SetComment(food.Comment).
			OnConflict().
			UpdateNewValues().
			ID(ctx)
	})

	return err
}

func (r *StorageSQLite) SetFoodComment(ctx context.Context, key, comment string) error {
	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		f, err := tx.Food.
			Query().
			Where(food.Key(key)).
			First(ctx)

		if err != nil {
			return nil, err
		}

		return f.
			Update().
			SetComment(comment).
			Save(ctx)
	})

	if ent.IsNotFound(err) {
		return ErrFoodNotFound
	}

	return err
}

func (r *StorageSQLite) GetFoodList(ctx context.Context) ([]Food, error) {
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Food.
			Query().
			Order(food.ByName()).
			All(ctx)
	})
	if err != nil {
		return nil, err
	}

	efList, _ := res.([]*ent.Food)

	if len(efList) == 0 {
		return nil, ErrFoodEmptyList
	}

	fList := make([]Food, 0, len(efList))
	for _, ef := range efList {
		fList = append(fList, *foodFromEntFood(ef))
	}

	return fList, nil
}

func (r *StorageSQLite) FindFood(ctx context.Context, pattern string) ([]Food, error) {
	upPattern := strings.ToUpper(pattern)

	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Food.
			Query().
			Where(func(s *entsql.Selector) {
				s.Where(
					entsql.ExprP(
						fmt.Sprintf("go_upper(%s) LIKE '%%' || ? || '%%'", s.C(food.FieldKey)),
						upPattern,
					)).
					Or().
					Where(
						entsql.ExprP(
							fmt.Sprintf("go_upper(%s) LIKE '%%' || ? || '%%'", s.C(food.FieldName)),
							upPattern,
						)).
					Or().
					Where(
						entsql.ExprP(
							fmt.Sprintf("go_upper(%s) LIKE '%%' || ? || '%%'", s.C(food.FieldBrand)),
							upPattern,
						)).
					Or().
					Where(
						entsql.ExprP(
							fmt.Sprintf("go_upper(%s) LIKE '%%' || ? || '%%'", s.C(food.FieldComment)),
							upPattern,
						))
			}).
			Order(food.ByName()).
			All(ctx)
	})
	if err != nil {
		return nil, err
	}

	efList, _ := res.([]*ent.Food)

	if len(efList) == 0 {
		return nil, ErrFoodEmptyList
	}

	fList := make([]Food, 0, len(efList))
	for _, ef := range efList {
		fList = append(fList, *foodFromEntFood(ef))
	}

	return fList, nil
}

func (r *StorageSQLite) DeleteFood(ctx context.Context, key string) error {
	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Food.
			Delete().
			Where(food.Key(key)).
			Exec(ctx)
	})

	// TODO: check constraint when delete with key in journal

	return err
}

func foodFromEntFood(ef *ent.Food) *Food {
	return &Food{
		Key:     ef.Key,
		Name:    ef.Name,
		Brand:   ef.Brand,
		Cal100:  ef.Cal100,
		Prot100: ef.Prot100,
		Fat100:  ef.Fat100,
		Carb100: ef.Carb100,
		Comment: ef.Comment,
	}
}

//
// Weight
//

func (r *StorageSQLite) GetWeightList(ctx context.Context, userID int64, from, to time.Time) ([]Weight, error) {
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
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

	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
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
	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
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

	_, err := r.dbSQL.ExecContext(ctx, _sqlSetJournal, userID, journal.Timestamp, journal.Meal, journal.FoodKey, journal.FoodWeight)
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
	_, err := r.dbSQL.ExecContext(ctx, _sqlDeleteJournal, userID, timestamp, meal, foodkey)
	return err
}

func (r *StorageSQLite) GetJournalReport(ctx context.Context, userID int64, from, to int64) ([]JournalReport, error) {
	rows, err := r.dbSQL.QueryContext(ctx, _sqlGetJournalReport, userID, from, to)
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
	rows, err := r.dbSQL.QueryContext(ctx, _sqlGetJournalStats, userID, from, to)
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
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.UserSettings.
			Query().
			Where(usersettings.Userid(userID)).
			First(ctx)
	})

	if err != nil {
		if ent.IsNotFound(err) {
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

	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
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

func (r *StorageSQLite) migrateWeight() error {
	// Get from old table
	ctx := context.Background()
	rows, err := r.dbSQL.QueryContext(ctx, `
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
	_, err = r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
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
	_, err = r.dbSQL.ExecContext(ctx, `
	DELETE FROM weight2;
	`)

	return err
}
