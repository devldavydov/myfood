package storage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageSQLiteTestSuite struct {
	suite.Suite

	stg    *StorageSQLite
	dbFile string
}

//
// Weight
//

func (r *StorageSQLiteTestSuite) TestGetWeightList() {
	r.Run("empty list", func() {
		_, err := r.stg.GetWeightList(context.TODO(), 1, 0, 10)
		r.ErrorIs(err, ErrWeightEmptyList)
	})

	r.Run("add data", func() {
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 1, Value: 1}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 2, Value: 2}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 3, Value: 3}))

		r.NoError(r.stg.SetWeight(context.TODO(), 2, &Weight{Timestamp: 4, Value: 4}))
	})

	r.Run("get list for different users", func() {
		lst, err := r.stg.GetWeightList(context.TODO(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: 1, Value: 1},
			{Timestamp: 2, Value: 2},
			{Timestamp: 3, Value: 3},
		}, lst)

		lst, err = r.stg.GetWeightList(context.TODO(), 2, 4, 4)
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: 4, Value: 4},
		}, lst)
	})

	r.Run("get limited list", func() {
		lst, err := r.stg.GetWeightList(context.TODO(), 1, 2, 3)
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: 2, Value: 2},
			{Timestamp: 3, Value: 3},
		}, lst)
	})
}

func (r *StorageSQLiteTestSuite) TestDeleteWeight() {
	r.Run("add data", func() {
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 1, Value: 1}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 2, Value: 2}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 3, Value: 3}))

		r.NoError(r.stg.SetWeight(context.TODO(), 2, &Weight{Timestamp: 4, Value: 4}))
	})

	r.Run("delete with incorrect user", func() {
		r.NoError(r.stg.DeleteWeight(context.TODO(), 2, 2))

		// Data not changed
		lst, err := r.stg.GetWeightList(context.TODO(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: 1, Value: 1},
			{Timestamp: 2, Value: 2},
			{Timestamp: 3, Value: 3},
		}, lst)

		lst, err = r.stg.GetWeightList(context.TODO(), 2, 4, 4)
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: 4, Value: 4},
		}, lst)
	})

	r.Run("delete weight for user", func() {
		r.NoError(r.stg.DeleteWeight(context.TODO(), 2, 4))
		_, err := r.stg.GetWeightList(context.TODO(), 2, 4, 4)
		r.ErrorIs(err, ErrWeightEmptyList)
	})
}

func (r *StorageSQLiteTestSuite) TestWeightCRU() {
	r.Run("get not existing weight", func() {
		wl, err := r.stg.GetWeightList(context.TODO(), 1, 1, 5)
		r.ErrorIs(err, ErrWeightEmptyList)
		r.Nil(wl)
	})

	r.Run("set invalid weight", func() {
		r.ErrorIs(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: -1, Value: 1}), ErrWeightInvalid)
		r.ErrorIs(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 1, Value: -1}), ErrWeightInvalid)
	})

	r.Run("set valid weight", func() {
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 1, Value: 1}))
	})

	r.Run("get weight", func() {
		wl, err := r.stg.GetWeightList(context.TODO(), 1, 1, 5)
		r.NoError(err)
		r.Equal([]Weight{{Timestamp: 1, Value: 1}}, wl)
	})

	r.Run("set again with update one", func() {
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 1, Value: 2}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: 2, Value: 2}))
	})

	r.Run("get weight list", func() {
		wl, err := r.stg.GetWeightList(context.TODO(), 1, 1, 5)
		r.NoError(err)
		r.Equal([]Weight{{Timestamp: 1, Value: 2}, {Timestamp: 2, Value: 2}}, wl)
	})
}

//
// UserSettings
//

func (r *StorageSQLiteTestSuite) TestUserSettingsCRU() {
	r.Run("get not exists settings", func() {
		stgs, err := r.stg.GetUserSettings(context.TODO(), 1)
		r.Nil(stgs)
		r.ErrorIs(err, ErrUserSettingsNotFound)
	})

	r.Run("set invalid settings", func() {
		r.ErrorIs(r.stg.SetUserSettings(context.TODO(), 1, &UserSettings{CalLimit: -1}), ErrUserSettingsInvalid)
	})

	r.Run("set valid settings and get", func() {
		r.NoError(r.stg.SetUserSettings(context.TODO(), 1, &UserSettings{CalLimit: 100}))

		stgs, err := r.stg.GetUserSettings(context.TODO(), 1)
		r.NoError(err)
		r.Equal(float64(100), stgs.CalLimit)
	})

	r.Run("update valid settings and get", func() {
		stgs, err := r.stg.GetUserSettings(context.TODO(), 1)
		r.NoError(err)
		r.Equal(float64(100), stgs.CalLimit)

		r.NoError(r.stg.SetUserSettings(context.TODO(), 1, &UserSettings{CalLimit: 200}))

		stgs, err = r.stg.GetUserSettings(context.TODO(), 1)
		r.NoError(err)
		r.Equal(float64(200), stgs.CalLimit)
	})
}

//
// Food
//

func (r *StorageSQLiteTestSuite) TestFoodCRU() {
	r.Run("get food not exists", func() {
		food, err := r.stg.GetFood(context.TODO(), "key1")
		r.Nil(food)
		r.ErrorIs(err, ErrFoodNotFound)
	})

	r.Run("create invalid food", func() {
		r.ErrorIs(r.stg.SetFood(context.TODO(), &Food{
			Key: "", Name: "Name", Cal100: 1, Prot100: 1, Fat100: 1, Carb100: 1,
		}), ErrFoodInvalid)
		r.ErrorIs(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key", Name: "", Cal100: 1, Prot100: 1, Fat100: 1, Carb100: 1,
		}), ErrFoodInvalid)
		r.ErrorIs(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key", Name: "Name", Cal100: -1, Prot100: 1, Fat100: 1, Carb100: 1,
		}), ErrFoodInvalid)
		r.ErrorIs(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key", Name: "Name", Cal100: 1, Prot100: -1, Fat100: 1, Carb100: 1,
		}), ErrFoodInvalid)
		r.ErrorIs(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key", Name: "Name", Cal100: 1, Prot100: 1, Fat100: -1, Carb100: 1,
		}), ErrFoodInvalid)
		r.ErrorIs(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key", Name: "Name", Cal100: 1, Prot100: 1, Fat100: 1, Carb100: -1,
		}), ErrFoodInvalid)
	})

	r.Run("create valid food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key", Name: "Name", Brand: "Brand", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment",
		}))
	})

	r.Run("get food", func() {
		food, err := r.stg.GetFood(context.TODO(), "Key")
		r.NoError(err)
		r.Equal(&Food{
			Key: "Key", Name: "Name", Brand: "Brand", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment",
		}, food)
	})

	r.Run("update food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key", Name: "Name 2", Brand: "Brand 2", Cal100: 10, Prot100: 20, Fat100: 30, Carb100: 40, Comment: "Comment 2",
		}))

		food, err := r.stg.GetFood(context.TODO(), "Key")
		r.NoError(err)
		r.Equal(&Food{
			Key: "Key", Name: "Name 2", Brand: "Brand 2", Cal100: 10, Prot100: 20, Fat100: 30, Carb100: 40, Comment: "Comment 2",
		}, food)
	})
}

func (r *StorageSQLiteTestSuite) TestFoodList() {
	r.Run("get empty list", func() {
		lst, err := r.stg.GetFoodList(context.TODO())
		r.ErrorIs(err, ErrFoodEmptyList)
		r.Nil(lst)
	})

	r.Run("add food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key1", Name: "bbb", Brand: "Brand1", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment1",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key2", Name: "aaa", Brand: "Brand2", Cal100: 4, Prot100: 5, Fat100: 6, Carb100: 7, Comment: "Comment2",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key3", Name: "ccc", Brand: "Brand3", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "Comment3",
		}))
	})

	r.Run("get list", func() {
		lst, err := r.stg.GetFoodList(context.TODO())
		r.NoError(err)
		r.Equal([]Food{
			{Key: "Key2", Name: "aaa", Brand: "Brand2", Cal100: 4, Prot100: 5, Fat100: 6, Carb100: 7, Comment: "Comment2"},
			{Key: "Key1", Name: "bbb", Brand: "Brand1", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment1"},
			{Key: "Key3", Name: "ccc", Brand: "Brand3", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "Comment3"},
		}, lst)
	})
}

func (r *StorageSQLiteTestSuite) TestFindFood() {
	r.Run("add food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "kFind", Name: "bbb", Brand: "Brand1", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment1",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key2", Name: "nfind", Brand: "Brand2", Cal100: 4, Prot100: 5, Fat100: 6, Carb100: 7, Comment: "Comment2",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key3", Name: "ccc", Brand: "bfind", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "Comment3",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key4", Name: "ddd", Brand: "Brand3", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "cfind",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "едрус", Name: "Еда Русская", Brand: "Brand3", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "руСКом",
		}))
	})

	r.Run("find by key", func() {
		lst, err := r.stg.FindFood(context.TODO(), "kfind")
		r.NoError(err)
		r.Equal([]Food{
			{Key: "kFind", Name: "bbb", Brand: "Brand1", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment1"},
		}, lst)
	})

	r.Run("find by name", func() {
		lst, err := r.stg.FindFood(context.TODO(), "Nfind")
		r.NoError(err)
		r.Equal([]Food{
			{Key: "Key2", Name: "nfind", Brand: "Brand2", Cal100: 4, Prot100: 5, Fat100: 6, Carb100: 7, Comment: "Comment2"},
		}, lst)
	})

	r.Run("find by brand", func() {
		lst, err := r.stg.FindFood(context.TODO(), "bfind")
		r.NoError(err)
		r.Equal([]Food{
			{Key: "Key3", Name: "ccc", Brand: "bfind", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "Comment3"},
		}, lst)
	})

	r.Run("find by comment", func() {
		lst, err := r.stg.FindFood(context.TODO(), "cfind")
		r.NoError(err)
		r.Equal([]Food{
			{Key: "Key4", Name: "ddd", Brand: "Brand3", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "cfind"},
		}, lst)
	})

	r.Run("find all k", func() {
		lst, err := r.stg.FindFood(context.TODO(), "k")
		r.NoError(err)
		r.Equal([]Food{
			{Key: "kFind", Name: "bbb", Brand: "Brand1", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment1"},
			{Key: "Key3", Name: "ccc", Brand: "bfind", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "Comment3"},
			{Key: "Key4", Name: "ddd", Brand: "Brand3", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "cfind"},
			{Key: "Key2", Name: "nfind", Brand: "Brand2", Cal100: 4, Prot100: 5, Fat100: 6, Carb100: 7, Comment: "Comment2"},
		}, lst)
	})

	r.Run("find non latin", func() {
		for _, pattern := range []string{"рус", "ЕДА", "еДа", "сК"} {
			lst, err := r.stg.FindFood(context.TODO(), pattern)
			r.NoError(err)
			r.Equal([]Food{
				{Key: "едрус", Name: "Еда Русская", Brand: "Brand3", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "руСКом"},
			}, lst)
		}
	})
}

func (r *StorageSQLiteTestSuite) TestDeleteFood() {
	r.Run("delete not exists food", func() {
		r.NoError(r.stg.DeleteFood(context.TODO(), "key"))
	})

	r.Run("add food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key1", Name: "bbb", Brand: "Brand1", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment1",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key2", Name: "aaa", Brand: "Brand2", Cal100: 4, Prot100: 5, Fat100: 6, Carb100: 7, Comment: "Comment2",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key3", Name: "ccc", Brand: "Brand3", Cal100: 8, Prot100: 9, Fat100: 10, Carb100: 11, Comment: "Comment3",
		}))
	})

	r.Run("delete food", func() {
		f, err := r.stg.GetFood(context.TODO(), "Key1")
		r.NoError(err)
		r.Equal("Key1", f.Key)

		r.NoError(r.stg.DeleteFood(context.TODO(), "Key1"))

		_, err = r.stg.GetFood(context.TODO(), "Key1")
		r.ErrorIs(err, ErrFoodNotFound)
	})
}

//
// Journal
//

func (r *StorageSQLiteTestSuite) TestJournalCRUD() {
	r.Run("set invalid journal", func() {
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: -1, Meal: Meal(0), FoodKey: "food", FoodWeight: 100,
		}), ErrJournalInvalid)
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 1, Meal: Meal(-1), FoodKey: "food", FoodWeight: 100,
		}), ErrJournalInvalid)
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 1, Meal: Meal(1), FoodKey: "", FoodWeight: 100,
		}), ErrJournalInvalid)
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 1, Meal: Meal(1), FoodKey: "food", FoodWeight: 0,
		}), ErrJournalInvalid)
	})

	r.Run("set journal with invalid food", func() {
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 1, Meal: Meal(0), FoodKey: "food", FoodWeight: 100,
		}), ErrJournalInvalidFood)
	})

	r.Run("add food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "food_a", Name: "aaa", Brand: "brand a", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "food_b", Name: "bbb", Brand: "brand b", Cal100: 5, Prot100: 6, Fat100: 7, Carb100: 8, Comment: "",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "food_c", Name: "ccc", Brand: "brand c", Cal100: 1, Prot100: 1, Fat100: 1, Carb100: 1, Comment: "ccc",
		}))
	})

	r.Run("set journal for different timestamps and users", func() {
		// user 1, timestamp 1
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 1, Meal: Meal(0), FoodKey: "food_b", FoodWeight: 1,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 1, Meal: Meal(1), FoodKey: "food_a", FoodWeight: 2,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 1, Meal: Meal(2), FoodKey: "food_c", FoodWeight: 3,
		}))

		// user 1, timestamp 2
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 2, Meal: Meal(0), FoodKey: "food_b", FoodWeight: 3,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 2, Meal: Meal(1), FoodKey: "food_a", FoodWeight: 2,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 2, Meal: Meal(1), FoodKey: "food_c", FoodWeight: 1,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 2, Meal: Meal(2), FoodKey: "food_c", FoodWeight: 4,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: 2, Meal: Meal(2), FoodKey: "food_a", FoodWeight: 5,
		}))

		// user 2, timestamp 3
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &Journal{
			Timestamp: 3, Meal: Meal(0), FoodKey: "food_b", FoodWeight: 3,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &Journal{
			Timestamp: 3, Meal: Meal(1), FoodKey: "food_a", FoodWeight: 2,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &Journal{
			Timestamp: 3, Meal: Meal(1), FoodKey: "food_c", FoodWeight: 1,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &Journal{
			Timestamp: 3, Meal: Meal(1), FoodKey: "food_b", FoodWeight: 4,
		}))
	})

	r.Run("get empty report", func() {
		_, err := r.stg.GetJournalForPeriodAndMeal(context.TODO(), 1, 10, 20, Meal(0))
		r.ErrorIs(err, ErrJournalReportEmpty)
		_, err = r.stg.GetJournalForPeriod(context.TODO(), 1, 10, 20)
		r.ErrorIs(err, ErrJournalReportEmpty)
	})

	r.Run("get journal reports for user 1", func() {
		// report for period and meal
		rep, err := r.stg.GetJournalForPeriodAndMeal(context.TODO(), 1, 1, 1, Meal(1))
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: 1, Meal: Meal(1), FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 2, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
		}, rep)

		rep, err = r.stg.GetJournalForPeriodAndMeal(context.TODO(), 1, 2, 2, Meal(2))
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: 2, Meal: Meal(2), FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 5, Cal: 5, Prot: 10, Fat: 15, Carb: 20},
			{Timestamp: 2, Meal: Meal(2), FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 4, Cal: 4, Prot: 4, Fat: 4, Carb: 4},
		}, rep)

		// report for period
		rep, err = r.stg.GetJournalForPeriod(context.TODO(), 1, 1, 2)
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: 1, Meal: Meal(0), FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 1, Cal: 5, Prot: 6, Fat: 7, Carb: 8},
			{Timestamp: 1, Meal: Meal(1), FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 2, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: 1, Meal: Meal(2), FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 3, Cal: 3, Prot: 3, Fat: 3, Carb: 3},
			{Timestamp: 2, Meal: Meal(0), FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 3, Cal: 15, Prot: 18, Fat: 21, Carb: 24},
			{Timestamp: 2, Meal: Meal(1), FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 2, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: 2, Meal: Meal(1), FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 1, Cal: 1, Prot: 1, Fat: 1, Carb: 1},
			{Timestamp: 2, Meal: Meal(2), FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 5, Cal: 5, Prot: 10, Fat: 15, Carb: 20},
			{Timestamp: 2, Meal: Meal(2), FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 4, Cal: 4, Prot: 4, Fat: 4, Carb: 4},
		}, rep)
	})

	r.Run("check that user 2 gets his data", func() {
		// report for
		rep, err := r.stg.GetJournalForPeriod(context.TODO(), 2, 1, 3)
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: 3, Meal: Meal(0), FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 3, Cal: 15, Prot: 18, Fat: 21, Carb: 24},
			{Timestamp: 3, Meal: Meal(1), FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 2, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: 3, Meal: Meal(1), FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 4, Cal: 20, Prot: 24, Fat: 28, Carb: 32},
			{Timestamp: 3, Meal: Meal(1), FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 1, Cal: 1, Prot: 1, Fat: 1, Carb: 1},
		}, rep)
	})

	r.Run("update and delete for user 1", func() {
		r.NoError(r.stg.DeleteJournal(context.TODO(), 1, 1, Meal(0), "food_b"))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{Timestamp: 1, Meal: Meal(1), FoodKey: "food_a", FoodWeight: 3}))

		rep, err := r.stg.GetJournalForPeriod(context.TODO(), 1, 1, 1)
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: 1, Meal: Meal(1), FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 3, Cal: 3, Prot: 6, Fat: 9, Carb: 12},
			{Timestamp: 1, Meal: Meal(2), FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 3, Cal: 3, Prot: 3, Fat: 3, Carb: 3},
		}, rep)
	})
}

//
// Suite setup
//

func (r *StorageSQLiteTestSuite) SetupTest() {
	var err error

	f, err := os.CreateTemp("", "testdb")
	require.NoError(r.T(), err)
	r.dbFile = f.Name()
	f.Close()

	r.stg, err = NewStorageSQLite(r.dbFile, nil)
	require.NoError(r.T(), err)
}

func (r *StorageSQLiteTestSuite) TearDownTest() {
	r.stg.Close()
	require.NoError(r.T(), os.Remove(r.dbFile))
}

func TestStorageSQLite(t *testing.T) {
	suite.Run(t, new(StorageSQLiteTestSuite))
}
