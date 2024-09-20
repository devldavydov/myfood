package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/devldavydov/myfood/internal/storage/ent"
	"github.com/devldavydov/myfood/internal/storage/ent/activity"
	"github.com/devldavydov/myfood/internal/storage/ent/bundle"
	"github.com/devldavydov/myfood/internal/storage/ent/food"
	"github.com/devldavydov/myfood/internal/storage/ent/journal"
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
	db    *ent.Client
	debug bool
}

var _ Storage = (*StorageSQLite)(nil)

func go_upper(str string) string {
	return strings.ToUpper(str)
}

func NewStorageSQLite(dbFilePath string, opts ...func(*StorageSQLite)) (*StorageSQLite, error) {
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
	for _, opt := range opts {
		opt(stg)
	}

	return stg, nil
}

func WithDebug() func(*StorageSQLite) {
	return func(s *StorageSQLite) {
		s.debug = true
	}
}

func (r *StorageSQLite) Close() error {
	if r.db == nil {
		return nil
	}

	return r.db.Close()
}

func (r *StorageSQLite) doTx(ctx context.Context, fn TxFn) (any, error) {
	// Begin database transaction.

	var clt *ent.Client
	if r.debug {
		clt = r.db.Debug()
	} else {
		clt = r.db
	}

	tx, err := clt.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("db start tx error: %w", err)
	}

	// Rollback on potential panics.
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
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
		// Check food in bundles.
		bndls, err := tx.Bundle.Query().All(ctx)
		if err != nil {
			return nil, err
		}

		for _, bndl := range bndls {
			for k, v := range bndl.Data {
				if v > 0 && k == key {
					return nil, ErrFoodIsUsed
				}
			}
		}

		// Delete food.
		return tx.Food.
			Delete().
			Where(food.Key(key)).
			Exec(ctx)
	})

	if ent.IsConstraintError(err) {
		return ErrFoodIsUsed
	}

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
// Bundle
//

func (r *StorageSQLite) SetBundle(ctx context.Context, userID int64, bndl *Bundle) error {
	if !bndl.Validate() {
		return ErrBundleInvalid
	}

	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		// Check bundle data
		for k, v := range bndl.Data {
			if v == 0 {
				// Dependent bundle.

				// Check in DB.
				_, err := tx.Bundle.
					Query().
					Where(
						bundle.Key(k),
						bundle.Userid(userID),
					).
					First(ctx)
				if err != nil {
					if ent.IsNotFound(err) {
						return nil, ErrBundleDepBundleNotFound
					}
					return nil, err
				}

				// If same key - recursive not allowed.
				if k == bndl.Key {
					return nil, ErrBundleDepRecursive
				}
			} else {
				// Dependent food, check in DB
				_, err := tx.Food.
					Query().
					Where(food.Key(k)).
					First(ctx)
				if err != nil {
					if ent.IsNotFound(err) {
						return nil, ErrBundleDepFoodNotFound
					}
					return nil, err
				}
			}
		}

		// Set bundle in DB
		return tx.Bundle.
			Create().
			SetUserid(userID).
			SetKey(bndl.Key).
			SetData(bndl.Data).
			OnConflict().
			UpdateNewValues().
			ID(ctx)
	})

	return err
}

func (r *StorageSQLite) GetBundle(ctx context.Context, userID int64, key string) (*Bundle, error) {
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return r.getBundle(ctx, tx, userID, key)
	})

	if err != nil {
		return nil, err
	}

	bndl, _ := res.(*ent.Bundle)

	return &Bundle{Key: bndl.Key, Data: bndl.Data}, nil
}

func (r *StorageSQLite) getBundle(ctx context.Context, tx *ent.Tx, userID int64, key string) (*ent.Bundle, error) {
	res, err := tx.Bundle.
		Query().
		Where(
			bundle.Userid(userID),
			bundle.Key(key),
		).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrBundleNotFound
		}
		return nil, err
	}

	return res, nil
}

func (r *StorageSQLite) GetBundleList(ctx context.Context, userID int64) ([]Bundle, error) {
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Bundle.
			Query().
			Where(bundle.Userid(userID)).
			Order(bundle.ByKey()).
			All(ctx)
	})
	if err != nil {
		return nil, err
	}

	beLst, _ := res.([]*ent.Bundle)
	if len(beLst) == 0 {
		return nil, ErrBundleEmptyList
	}

	bLst := make([]Bundle, 0, len(beLst))
	for _, b := range beLst {
		bLst = append(bLst, Bundle{Key: b.Key, Data: b.Data})
	}

	return bLst, nil
}

func (r *StorageSQLite) getBundleFoodItems(ctx context.Context, tx *ent.Tx, userID int64, bndlData, foodItems map[string]float64) error {
	for k, v := range bndlData {
		if v > 0 {
			// If food - add to result map.
			foodItems[k] = v
			continue
		}

		// If bundle - get from DB and call recursive.
		bndl, err := r.getBundle(ctx, tx, userID, k)
		if err != nil {
			return err
		}

		if err := r.getBundleFoodItems(ctx, tx, userID, bndl.Data, foodItems); err != nil {
			return err
		}
	}

	return nil
}

func (r *StorageSQLite) DeleteBundle(ctx context.Context, userID int64, key string) error {
	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		// Check that bundle not used in other bundles.
		bndls, err := tx.Bundle.
			Query().
			Where(bundle.Userid(userID)).
			All(ctx)
		if err != nil {
			return nil, err
		}

		for _, bndl := range bndls {
			for k, v := range bndl.Data {
				if v == 0 && k == key {
					return nil, ErrBundleIsUsed
				}
			}
		}

		return tx.Bundle.
			Delete().
			Where(
				bundle.Userid(userID),
				bundle.Key(key),
			).Exec(ctx)
	})

	return err
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

	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		food, err := r.getFoodForJournal(ctx, tx, journal.FoodKey)
		if err != nil {
			return nil, err
		}

		return tx.Journal.
			Create().
			SetUserid(userID).
			SetTimestamp(journal.Timestamp).
			SetMeal(int64(journal.Meal)).
			SetFoodweight(journal.FoodWeight).
			SetFood(food).
			OnConflict().
			UpdateNewValues().
			ID(ctx)
	})

	return err
}

func (r *StorageSQLite) SetJournalBundle(ctx context.Context, userID int64, timestamp time.Time, meal Meal, bndlKey string) error {
	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		// Get bundle
		bndl, err := r.getBundle(ctx, tx, userID, bndlKey)
		if err != nil {
			return nil, err
		}

		// Get food items from bundle.
		resFood := make(map[string]float64)
		err = r.getBundleFoodItems(ctx, tx, userID, bndl.Data, resFood)
		if err != nil {
			return nil, err
		}

		// Add food to journal.
		for foodKey, foodWeight := range resFood {
			food, err := r.getFoodForJournal(ctx, tx, foodKey)
			if err != nil {
				return nil, err
			}

			_, err = tx.Journal.
				Create().
				SetUserid(userID).
				SetTimestamp(timestamp).
				SetMeal(int64(meal)).
				SetFoodweight(foodWeight).
				SetFood(food).
				OnConflict().
				UpdateNewValues().
				ID(ctx)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}

func (r *StorageSQLite) getFoodForJournal(ctx context.Context, tx *ent.Tx, key string) (*ent.Food, error) {
	food, err := tx.Food.
		Query().
		Where(food.Key(key)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrJournalInvalidFood
		}
		return nil, err
	}

	return food, nil
}

func (r *StorageSQLite) DeleteJournal(ctx context.Context, userID int64, timestamp time.Time, meal Meal, foodkey string) error {
	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Journal.
			Delete().
			Where(
				journal.Userid(userID),
				journal.Timestamp(timestamp),
				journal.Meal(int64(meal)),
				journal.HasFoodWith(food.Key(foodkey)),
			).
			Exec(ctx)
	})

	return err
}

func (r *StorageSQLite) DeleteJournalMeal(ctx context.Context, userID int64, timestamp time.Time, meal Meal) error {
	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Journal.
			Delete().
			Where(
				journal.Userid(userID),
				journal.Timestamp(timestamp),
				journal.Meal(int64(meal)),
			).
			Exec(ctx)
	})

	return err
}

func (r *StorageSQLite) GetJournalMealReport(ctx context.Context, userID int64, timestamp time.Time, meal Meal) (*JournalMealReport, error) {
	// Get meal items
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Journal.
			Query().
			Where(
				journal.Userid(userID),
				journal.Timestamp(timestamp),
				journal.Meal(int64(meal)),
			).
			Order(
				journal.ByFoodField(food.FieldName),
			).
			WithFood().
			All(ctx)
	})

	if err != nil {
		return nil, err
	}

	jLst, _ := res.([]*ent.Journal)
	if len(jLst) == 0 {
		return nil, ErrJournalMealReportEmpty
	}

	lst := make([]JournalMealItem, 0, len(jLst))
	var mealCal float64
	for _, item := range jLst {
		cal := item.Foodweight / 100 * item.Edges.Food.Cal100
		mealCal += cal
		lst = append(lst, JournalMealItem{
			Timestamp:  item.Timestamp,
			FoodKey:    item.Edges.Food.Key,
			FoodName:   item.Edges.Food.Name,
			FoodBrand:  item.Edges.Food.Brand,
			FoodWeight: item.Foodweight,
			Cal:        cal,
		})
	}

	// Get total consumed calories for day
	total, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Journal.
			Query().
			Where(
				journal.Userid(userID),
				journal.Timestamp(timestamp),
			).
			Modify(func(s *entsql.Selector) {
				f := entsql.Table(food.Table)
				s.
					Join(f).
					On(
						s.C(journal.FoodColumn),
						f.C(food.FieldID),
					).
					Select(
						entsql.Sum(
							fmt.Sprintf("%s / 100 * %s", s.C(journal.FieldFoodweight), f.C(food.FieldCal100)),
						),
					)
			}).
			Float64(ctx)
	})
	if err != nil {
		return nil, err
	}

	return &JournalMealReport{
		Items:           lst,
		ConsumedMealCal: mealCal,
		ConsumedDayCal:  total.(float64),
	}, nil
}

func (r *StorageSQLite) GetJournalReport(ctx context.Context, userID int64, from, to time.Time) ([]JournalReport, error) {
	var res []struct {
		ent.Journal
		FoodKey   string  `sql:"foodkey"`
		FoodName  string  `sql:"foodname"`
		FoodBrand string  `sql:"foodbrand"`
		Cal       float64 `sql:"cal"`
		Prot      float64 `sql:"prot"`
		Fat       float64 `sql:"fat"`
		Carb      float64 `sql:"carb"`
	}

	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		err := tx.Journal.
			Query().
			Where(
				journal.Userid(userID),
				journal.TimestampGTE(from),
				journal.TimestampLTE(to),
			).
			Modify(func(s *entsql.Selector) {
				f := entsql.Table(food.Table)
				s.
					Join(f).
					On(
						s.C(journal.FoodColumn),
						f.C(food.FieldID),
					).
					AppendSelect(
						entsql.As(f.C(food.FieldKey), "foodkey"),
						entsql.As(f.C(food.FieldName), "foodname"),
						entsql.As(f.C(food.FieldBrand), "foodbrand"),
						entsql.As(
							fmt.Sprintf("%s / 100 * %s", s.C(journal.FieldFoodweight), f.C(food.FieldCal100)),
							"cal",
						),
						entsql.As(
							fmt.Sprintf("%s / 100 * %s", s.C(journal.FieldFoodweight), f.C(food.FieldProt100)),
							"prot",
						),
						entsql.As(
							fmt.Sprintf("%s / 100 * %s", s.C(journal.FieldFoodweight), f.C(food.FieldFat100)),
							"fat",
						),
						entsql.As(
							fmt.Sprintf("%s / 100 * %s", s.C(journal.FieldFoodweight), f.C(food.FieldCarb100)),
							"carb",
						),
					).
					OrderBy(
						s.C(journal.FieldTimestamp),
						s.C(journal.FieldMeal),
						f.C(food.FieldName),
					)
			}).
			Scan(ctx, &res)

		return nil, err
	})

	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrJournalReportEmpty
	}

	lst := make([]JournalReport, 0, len(res))
	for _, item := range res {
		lst = append(lst, JournalReport{
			Timestamp:  item.Timestamp,
			Meal:       Meal(item.Meal),
			FoodKey:    item.FoodKey,
			FoodName:   item.FoodName,
			FoodBrand:  item.FoodBrand,
			FoodWeight: item.Foodweight,
			Cal:        item.Cal,
			Prot:       item.Prot,
			Fat:        item.Fat,
			Carb:       item.Carb,
		})
	}

	return lst, nil
}

func (r *StorageSQLite) GetJournalStats(ctx context.Context, userID int64, from, to time.Time) ([]JournalStats, error) {
	var res []struct {
		Timestamp time.Time `sql:"timestamp"`
		TotalCal  float64   `sql:"totalCal"`
		TotalProt float64   `sql:"totalProt"`
		TotalFat  float64   `sql:"totalFat"`
		TotalCarb float64   `sql:"totalCarb"`
	}

	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		err := tx.Journal.
			Query().
			Where(
				journal.Userid(userID),
				journal.TimestampGTE(from),
				journal.TimestampLTE(to),
			).
			Modify(func(s *entsql.Selector) {
				f := entsql.Table(food.Table)
				s.
					Join(f).
					On(
						s.C(journal.FoodColumn),
						f.C(food.FieldID),
					).
					Select(
						entsql.As(s.C(journal.FieldTimestamp), "timestamp"),
						entsql.As(
							entsql.Sum(
								fmt.Sprintf("%s / 100 * %s", s.C(journal.FieldFoodweight), f.C(food.FieldCal100)),
							),
							"totalCal",
						),
						entsql.As(
							entsql.Sum(
								fmt.Sprintf("%s / 100 * %s", s.C(journal.FieldFoodweight), f.C(food.FieldProt100)),
							),
							"totalProt",
						),
						entsql.As(
							entsql.Sum(
								fmt.Sprintf("%s / 100 * %s", s.C(journal.FieldFoodweight), f.C(food.FieldFat100)),
							),
							"totalFat",
						),
						entsql.As(
							entsql.Sum(
								fmt.Sprintf("%s / 100 * %s", s.C(journal.FieldFoodweight), f.C(food.FieldCarb100)),
							),
							"totalCarb",
						),
					).
					GroupBy(
						s.C(journal.FieldTimestamp),
					).
					OrderBy(
						s.C(journal.FieldTimestamp),
					)
			}).
			Scan(ctx, &res)

		return nil, err
	})

	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrJournalStatsEmpty
	}

	lst := make([]JournalStats, 0, len(res))
	for _, item := range res {
		lst = append(lst, JournalStats{
			Timestamp: item.Timestamp,
			TotalCal:  item.TotalCal,
			TotalProt: item.TotalProt,
			TotalFat:  item.TotalFat,
			TotalCarb: item.TotalCarb,
		})
	}

	return lst, nil
}

func (r *StorageSQLite) CopyJournal(ctx context.Context, userID int64, from time.Time, mealFrom Meal, to time.Time, mealTo Meal) (int, error) {
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		// Check that destination is empty
		cnt, err := tx.Journal.
			Query().
			Where(
				journal.Userid(userID),
				journal.Timestamp(to),
				journal.Meal(int64(mealTo)),
			).
			Count(ctx)

		if err != nil {
			return nil, err
		}

		if cnt != 0 {
			return nil, ErrCopyToNotEmpty
		}

		// Get source list
		lst, err := tx.Journal.
			Query().
			Where(
				journal.Userid(userID),
				journal.Timestamp(from),
				journal.Meal(int64(mealFrom)),
			).
			WithFood().
			All(ctx)

		if err != nil {
			return 0, err
		}

		cnt = len(lst)
		bulk := make([]*ent.JournalCreate, 0, cnt)
		for _, item := range lst {
			bulk = append(bulk, tx.Journal.
				Create().
				SetUserid(userID).
				SetTimestamp(to).
				SetMeal(int64(mealTo)).
				SetFoodweight(item.Foodweight).
				SetFoodID(item.Edges.Food.ID),
			)
		}
		_, err = tx.Journal.CreateBulk(bulk...).Save(ctx)

		return cnt, err
	})

	cnt := 0
	if err == nil {
		cnt, _ = res.(int)
	}
	return cnt, err
}

func (r *StorageSQLite) GetJournalFoodAvgWeight(ctx context.Context, userID int64, from, to time.Time, foodkey string) (float64, error) {
	var v []struct {
		Avg float64
	}
	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		// Get food
		foodObj, err := tx.Food.
			Query().
			Where(food.Key(foodkey)).
			First(ctx)
		if err != nil {
			return nil, err
		}

		// Get avg food weight
		err = tx.Journal.
			Query().
			Where(
				journal.Userid(userID),
				journal.HasFoodWith(food.ID(foodObj.ID)),
				journal.TimestampGTE(from),
				journal.TimestampLTE(to),
			).
			Aggregate(ent.Mean(journal.FieldFoodweight)).
			Scan(ctx, &v)

		return nil, err
	})

	if err != nil {
		if ent.IsNotFound(err) {
			return 0, ErrJournalInvalidFood
		}

		return 0, err
	}

	return v[0].Avg, nil
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

	return &UserSettings{
		CalLimit:         us.CalLimit,
		DefaultActiveCal: us.DefaultActiveCal,
	}, nil
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
			SetDefaultActiveCal(settings.DefaultActiveCal).
			OnConflict().
			UpdateNewValues().
			ID(ctx)
	})

	return err
}

//
// Activity.
//

func (r *StorageSQLite) GetActivityList(ctx context.Context, userID int64, from, to time.Time) ([]Activity, error) {
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Activity.
			Query().
			Where(
				activity.Userid(userID),
				activity.TimestampGTE(from),
				activity.TimestampLTE(to),
			).
			Order(
				activity.ByTimestamp(),
			).
			All(ctx)
	})
	if err != nil {
		return nil, err
	}

	aeLst, _ := res.([]*ent.Activity)
	if len(aeLst) == 0 {
		return nil, ErrActivityEmptyList
	}

	aLst := make([]Activity, 0, len(aeLst))
	for _, a := range aeLst {
		aLst = append(aLst, Activity{Timestamp: a.Timestamp, ActiveCal: a.ActiveCal})
	}

	return aLst, nil
}

func (r *StorageSQLite) GetActivity(ctx context.Context, userID int64, timestamp time.Time) (*Activity, error) {
	res, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Activity.
			Query().
			Where(activity.UseridEQ(userID), activity.Timestamp(timestamp)).
			First(ctx)
	})

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}

	a, _ := res.(*ent.Activity)
	return &Activity{
		Timestamp: a.Timestamp,
		ActiveCal: a.ActiveCal,
	}, nil
}

func (r *StorageSQLite) SetActivity(ctx context.Context, userID int64, activity *Activity) error {
	if !activity.Validate() {
		return ErrActivityInvalid
	}

	_, err := r.doTx(ctx, func(ctx context.Context, tx *ent.Tx) (any, error) {
		return tx.Activity.
			Create().
			SetUserid(userID).
			SetTimestamp(activity.Timestamp).
			SetActiveCal(activity.ActiveCal).
			OnConflict().
			UpdateNewValues().
			ID(ctx)
	})

	return err
}

func (r *StorageSQLite) DeleteActivity(ctx context.Context, userID int64, timestamp time.Time) error {
	return nil
}
