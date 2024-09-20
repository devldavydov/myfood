package storage

import (
	"context"
	"os"
	"testing"
	"time"

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
		_, err := r.stg.GetWeightList(context.TODO(), 1, T(0), T(10))
		r.ErrorIs(err, ErrWeightEmptyList)
	})

	r.Run("add data", func() {
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(1), Value: 1}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(2), Value: 2}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(3), Value: 3}))

		r.NoError(r.stg.SetWeight(context.TODO(), 2, &Weight{Timestamp: T(4), Value: 4}))
	})

	r.Run("get list for different users", func() {
		lst, err := r.stg.GetWeightList(context.TODO(), 1, T(1), T(3))
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: T(1), Value: 1},
			{Timestamp: T(2), Value: 2},
			{Timestamp: T(3), Value: 3},
		}, lst)

		lst, err = r.stg.GetWeightList(context.TODO(), 2, T(4), T(4))
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: T(4), Value: 4},
		}, lst)
	})

	r.Run("get limited list", func() {
		lst, err := r.stg.GetWeightList(context.TODO(), 1, T(2), T(3))
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: T(2), Value: 2},
			{Timestamp: T(3), Value: 3},
		}, lst)
	})
}

func (r *StorageSQLiteTestSuite) TestDeleteWeight() {
	r.Run("add data", func() {
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(1), Value: 1}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(2), Value: 2}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(3), Value: 3}))

		r.NoError(r.stg.SetWeight(context.TODO(), 2, &Weight{Timestamp: T(4), Value: 4}))
	})

	r.Run("delete with incorrect user", func() {
		r.NoError(r.stg.DeleteWeight(context.TODO(), 2, T(2)))

		// Data not changed
		lst, err := r.stg.GetWeightList(context.TODO(), 1, T(1), T(3))
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: T(1), Value: 1},
			{Timestamp: T(2), Value: 2},
			{Timestamp: T(3), Value: 3},
		}, lst)

		lst, err = r.stg.GetWeightList(context.TODO(), 2, T(4), T(4))
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: T(4), Value: 4},
		}, lst)
	})

	r.Run("delete weight for user", func() {
		r.NoError(r.stg.DeleteWeight(context.TODO(), 2, T(4)))
		_, err := r.stg.GetWeightList(context.TODO(), 2, T(4), T(4))
		r.ErrorIs(err, ErrWeightEmptyList)
	})
}

func (r *StorageSQLiteTestSuite) TestWeightCRU() {
	r.Run("get not existing weight", func() {
		wl, err := r.stg.GetWeightList(context.TODO(), 1, T(1), T(5))
		r.ErrorIs(err, ErrWeightEmptyList)
		r.Nil(wl)
	})

	r.Run("set invalid weight", func() {
		r.ErrorIs(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(1), Value: -1}), ErrWeightInvalid)
	})

	r.Run("set valid weight", func() {
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(1), Value: 1}))
	})

	r.Run("get weight", func() {
		wl, err := r.stg.GetWeightList(context.TODO(), 1, T(1), T(5))
		r.NoError(err)
		r.Equal([]Weight{{Timestamp: T(1), Value: 1}}, wl)
	})

	r.Run("set again with update one", func() {
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(1), Value: 2}))
		r.NoError(r.stg.SetWeight(context.TODO(), 1, &Weight{Timestamp: T(2), Value: 2}))
	})

	r.Run("get weight list", func() {
		wl, err := r.stg.GetWeightList(context.TODO(), 1, T(1), T(5))
		r.NoError(err)
		r.Equal([]Weight{{Timestamp: T(1), Value: 2}, {Timestamp: T(2), Value: 2}}, wl)
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
		r.ErrorIs(r.stg.SetUserSettings(
			context.TODO(),
			1,
			&UserSettings{CalLimit: -1, DefaultActiveCal: -1}),
			ErrUserSettingsInvalid)
	})

	r.Run("set valid settings and get", func() {
		r.NoError(r.stg.SetUserSettings(
			context.TODO(),
			1,
			&UserSettings{CalLimit: 100, DefaultActiveCal: 200}))

		stgs, err := r.stg.GetUserSettings(context.TODO(), 1)
		r.NoError(err)
		r.Equal(float64(100), stgs.CalLimit)
		r.Equal(float64(200), stgs.DefaultActiveCal)
	})

	r.Run("update valid settings and get", func() {
		stgs, err := r.stg.GetUserSettings(context.TODO(), 1)
		r.NoError(err)
		r.Equal(float64(100), stgs.CalLimit)
		r.Equal(float64(200), stgs.DefaultActiveCal)

		r.NoError(r.stg.SetUserSettings(
			context.TODO(),
			1,
			&UserSettings{CalLimit: 300, DefaultActiveCal: 300}))

		stgs, err = r.stg.GetUserSettings(context.TODO(), 1)
		r.NoError(err)
		r.Equal(float64(300), stgs.CalLimit)
		r.Equal(float64(300), stgs.DefaultActiveCal)
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

func (r *StorageSQLiteTestSuite) TestFoodSetComment() {
	r.Run("set comment for not exists food", func() {
		r.ErrorIs(r.stg.SetFoodComment(context.TODO(), "key", "comment"), ErrFoodNotFound)
	})

	r.Run("add food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "Key1", Name: "bbb", Brand: "Brand1", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "",
		}))
	})

	r.Run("check and set comment", func() {
		f, err := r.stg.GetFood(context.TODO(), "Key1")
		r.NoError(err)
		r.Equal("", f.Comment)

		r.NoError(r.stg.SetFoodComment(context.TODO(), "Key1", "FooBar"))

		f, err = r.stg.GetFood(context.TODO(), "Key1")
		r.NoError(err)
		r.Equal("FooBar", f.Comment)
	})
}

//
// Journal
//

func (r *StorageSQLiteTestSuite) TestJournalCRUD() {
	r.Run("set invalid journal", func() {
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(1), Meal: Meal(-1), FoodKey: "food", FoodWeight: 100,
		}), ErrJournalInvalid)
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(1), Meal: Meal(1), FoodKey: "", FoodWeight: 100,
		}), ErrJournalInvalid)
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(1), Meal: Meal(1), FoodKey: "food", FoodWeight: 0,
		}), ErrJournalInvalid)
	})

	r.Run("set journal with invalid food", func() {
		r.ErrorIs(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(1), Meal: Meal(0), FoodKey: "food", FoodWeight: 100,
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
			Timestamp: T(1), Meal: Meal(0), FoodKey: "food_b", FoodWeight: 100,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(1), Meal: Meal(1), FoodKey: "food_a", FoodWeight: 200,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(1), Meal: Meal(2), FoodKey: "food_c", FoodWeight: 300,
		}))

		// user 1, timestamp 2
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(2), Meal: Meal(0), FoodKey: "food_b", FoodWeight: 300,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(2), Meal: Meal(1), FoodKey: "food_a", FoodWeight: 200,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(2), Meal: Meal(1), FoodKey: "food_c", FoodWeight: 100,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(2), Meal: Meal(2), FoodKey: "food_c", FoodWeight: 400,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(2), Meal: Meal(2), FoodKey: "food_a", FoodWeight: 500,
		}))

		// user 2, timestamp 3
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &Journal{
			Timestamp: T(3), Meal: Meal(0), FoodKey: "food_b", FoodWeight: 300,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &Journal{
			Timestamp: T(3), Meal: Meal(1), FoodKey: "food_a", FoodWeight: 200,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &Journal{
			Timestamp: T(3), Meal: Meal(1), FoodKey: "food_c", FoodWeight: 100,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 2, &Journal{
			Timestamp: T(3), Meal: Meal(1), FoodKey: "food_b", FoodWeight: 400,
		}))
	})

	r.Run("get empty report", func() {
		_, err := r.stg.GetJournalReport(context.TODO(), 1, T(10), T(20))
		r.ErrorIs(err, ErrJournalReportEmpty)

		_, err = r.stg.GetJournalStats(context.TODO(), 1, T(10), T(20))
		r.ErrorIs(err, ErrJournalStatsEmpty)

		_, err = r.stg.GetJournalMealReport(context.TODO(), 1, T(10), Meal(0))
		r.ErrorIs(err, ErrJournalMealReportEmpty)
	})

	r.Run("get journal reports for user 1", func() {
		rep, err := r.stg.GetJournalReport(context.TODO(), 1, T(1), T(2))
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: T(1), Meal: Meal(0), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 100, Cal: 5, Prot: 6, Fat: 7, Carb: 8},
			{Timestamp: T(1), Meal: Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: T(1), Meal: Meal(2), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 3, Prot: 3, Fat: 3, Carb: 3},
			{Timestamp: T(2), Meal: Meal(0), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 300, Cal: 15, Prot: 18, Fat: 21, Carb: 24},
			{Timestamp: T(2), Meal: Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: T(2), Meal: Meal(1), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 100, Cal: 1, Prot: 1, Fat: 1, Carb: 1},
			{Timestamp: T(2), Meal: Meal(2), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 500, Cal: 5, Prot: 10, Fat: 15, Carb: 20},
			{Timestamp: T(2), Meal: Meal(2), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 400, Cal: 4, Prot: 4, Fat: 4, Carb: 4},
		}, rep)

		stats, err := r.stg.GetJournalStats(context.TODO(), 1, T(1), T(2))
		r.NoError(err)
		r.Equal([]JournalStats{
			{Timestamp: T(1), TotalCal: 10, TotalProt: 13, TotalFat: 16, TotalCarb: 19},
			{Timestamp: T(2), TotalCal: 27, TotalProt: 37, TotalFat: 47, TotalCarb: 57},
		}, stats)

		mealRep, err := r.stg.GetJournalMealReport(context.TODO(), 1, T(2), Meal(1))
		r.NoError(err)
		r.Equal([]JournalMealItem{
			{Timestamp: T(2), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2},
			{Timestamp: T(2), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 100, Cal: 1},
		}, mealRep.Items)
		r.Equal(float64(3), mealRep.ConsumedMealCal)
		r.Equal(float64(27), mealRep.ConsumedDayCal)

		foodAvgW, err := r.stg.GetJournalFoodAvgWeight(context.TODO(), 1, T(1), T(2), "food_b")
		r.NoError(err)
		r.Equal(float64(200), foodAvgW)
	})

	r.Run("check that user 2 gets his data", func() {
		// report for
		rep, err := r.stg.GetJournalReport(context.TODO(), 2, T(1), T(3))
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: T(3), Meal: Meal(0), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 300, Cal: 15, Prot: 18, Fat: 21, Carb: 24},
			{Timestamp: T(3), Meal: Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: T(3), Meal: Meal(1), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 400, Cal: 20, Prot: 24, Fat: 28, Carb: 32},
			{Timestamp: T(3), Meal: Meal(1), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 100, Cal: 1, Prot: 1, Fat: 1, Carb: 1},
		}, rep)
	})

	r.Run("update and delete for user 1", func() {
		r.NoError(r.stg.DeleteJournal(context.TODO(), 1, T(1), Meal(0), "food_b"))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{Timestamp: T(1), Meal: Meal(1), FoodKey: "food_a", FoodWeight: 300}))

		rep, err := r.stg.GetJournalReport(context.TODO(), 1, T(1), T(1))
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: T(1), Meal: Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 300, Cal: 3, Prot: 6, Fat: 9, Carb: 12},
			{Timestamp: T(1), Meal: Meal(2), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 3, Prot: 3, Fat: 3, Carb: 3},
		}, rep)
	})

	r.Run("try delete used food", func() {
		r.ErrorIs(r.stg.DeleteFood(context.TODO(), "food_a"), ErrFoodIsUsed)
	})

	r.Run("delete meal for day", func() {
		mealRep, err := r.stg.GetJournalMealReport(context.TODO(), 1, T(2), Meal(1))
		r.NoError(err)
		r.Equal(2, len(mealRep.Items))

		r.NoError(r.stg.DeleteJournalMeal(context.TODO(), 1, T(2), Meal(1)))

		_, err = r.stg.GetJournalMealReport(context.TODO(), 1, T(2), Meal(1))
		r.ErrorIs(err, ErrJournalMealReportEmpty)
	})
}

func (r *StorageSQLiteTestSuite) TestJournalCopy() {
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

	r.Run("set initial journal", func() {
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(1), Meal: Meal(0), FoodKey: "food_b", FoodWeight: 100,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(1), Meal: Meal(0), FoodKey: "food_a", FoodWeight: 200,
		}))
		r.NoError(r.stg.SetJournal(context.TODO(), 1, &Journal{
			Timestamp: T(2), Meal: Meal(0), FoodKey: "food_c", FoodWeight: 300,
		}))
	})

	r.Run("try copy when dest no empty", func() {
		_, err := r.stg.CopyJournal(context.TODO(), 1, T(1), Meal(0), T(2), Meal(0))
		r.ErrorIs(err, ErrCopyToNotEmpty)
	})

	r.Run("copy success", func() {
		cnt, err := r.stg.CopyJournal(context.TODO(), 1, T(1), Meal(0), T(2), Meal(1))
		r.NoError(err)
		r.Equal(2, cnt)

		rep, err := r.stg.GetJournalReport(context.TODO(), 1, T(2), T(2))
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: T(2), Meal: Meal(0), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 3, Prot: 3, Fat: 3, Carb: 3},
			{Timestamp: T(2), Meal: Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 200, Cal: 2, Prot: 4, Fat: 6, Carb: 8},
			{Timestamp: T(2), Meal: Meal(1), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 100, Cal: 5, Prot: 6, Fat: 7, Carb: 8},
		}, rep)
	})

	r.Run("copy zero", func() {
		cnt, err := r.stg.CopyJournal(context.TODO(), 1, T(10), Meal(0), T(20), Meal(1))
		r.NoError(err)
		r.Equal(0, cnt)
	})
}

//
// Bundle
//

func (r *StorageSQLiteTestSuite) TestBundleCRUD() {
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

	r.Run("set invalid bundle", func() {
		r.ErrorIs(r.stg.SetBundle(
			context.TODO(), 1, &Bundle{},
		), ErrBundleInvalid)
		r.ErrorIs(r.stg.SetBundle(
			context.TODO(), 1, &Bundle{Key: "bndl"},
		), ErrBundleInvalid)
		r.ErrorIs(r.stg.SetBundle(
			context.TODO(), 1, &Bundle{Data: map[string]float64{"food_a": 10}},
		), ErrBundleInvalid)
		r.ErrorIs(r.stg.SetBundle(
			context.TODO(), 1, &Bundle{Data: map[string]float64{"food_a": -10}, Key: "bndl"},
		), ErrBundleInvalid)
	})

	r.Run("set bundle with unknown food", func() {
		r.ErrorIs(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl1",
			Data: map[string]float64{
				"food_a": 10,
				"food_d": 11,
			},
		}), ErrBundleDepFoodNotFound)
	})

	r.Run("set bundle success", func() {
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl1",
			Data: map[string]float64{
				"food_a": 10,
				"food_b": 11,
			},
		}))
	})

	r.Run("set bundle success for another user", func() {
		r.NoError(r.stg.SetBundle(context.TODO(), 2, &Bundle{
			Key: "bndl3",
			Data: map[string]float64{
				"food_a": 10,
				"food_b": 11,
			},
		}))
	})

	r.Run("set bundle with not found dependent bundle", func() {
		r.ErrorIs(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl2",
			Data: map[string]float64{
				"food_a": 10,
				"food_b": 11,
				"bndl3":  0, // not work, because bundle3 for user 2.
			},
		}), ErrBundleDepBundleNotFound)
	})

	r.Run("set bundle with dependent bundle", func() {
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl2",
			Data: map[string]float64{
				"food_a": 10,
				"food_b": 11,
				"bndl1":  0,
			},
		}))
	})

	r.Run("set bundle with recursive dependent bundle", func() {
		r.ErrorIs(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl2",
			Data: map[string]float64{
				"food_a": 10,
				"food_b": 11,
				"bndl2":  0,
			},
		}), ErrBundleDepRecursive)
	})

	r.Run("get unknown bundle", func() {
		_, err := r.stg.GetBundle(context.TODO(), 1, "bndl3")
		r.ErrorIs(err, ErrBundleNotFound)
	})

	r.Run("get bundle", func() {
		bndl, err := r.stg.GetBundle(context.TODO(), 1, "bndl2")
		r.NoError(err)
		r.Equal(bndl, &Bundle{
			Key: "bndl2",
			Data: map[string]float64{
				"bndl1":  0,
				"food_a": 10,
				"food_b": 11,
			},
		})
	})

	r.Run("get bundle list", func() {
		bndls, err := r.stg.GetBundleList(context.TODO(), 1)
		r.NoError(err)
		r.Equal(bndls, []Bundle{
			{Key: "bndl1", Data: map[string]float64{
				"food_a": 10,
				"food_b": 11,
			}},
			{Key: "bndl2", Data: map[string]float64{
				"food_a": 10,
				"food_b": 11,
				"bndl1":  0,
			}},
		})
	})

	r.Run("update bundle", func() {
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl2",
			Data: map[string]float64{
				"food_a": 20,
				"food_b": 30,
				"bndl1":  0,
			},
		}))
	})

	r.Run("get bundle", func() {
		bndl, err := r.stg.GetBundle(context.TODO(), 1, "bndl2")
		r.NoError(err)
		r.Equal(bndl, &Bundle{
			Key: "bndl2",
			Data: map[string]float64{
				"bndl1":  0,
				"food_a": 20,
				"food_b": 30,
			},
		})
	})

	r.Run("try delete when bundle is used", func() {
		r.ErrorIs(r.stg.DeleteBundle(context.TODO(), 1, "bndl1"), ErrBundleIsUsed)
	})

	r.Run("try delete food used in bundle", func() {
		r.ErrorIs(r.stg.DeleteFood(context.TODO(), "food_a"), ErrFoodIsUsed)
	})

	r.Run("delete bundles success", func() {
		r.NoError(r.stg.DeleteBundle(context.TODO(), 1, "bndl2"))
		r.NoError(r.stg.DeleteBundle(context.TODO(), 1, "bndl1"))
	})

	r.Run("empty bundle list", func() {
		_, err := r.stg.GetBundleList(context.TODO(), 1)
		r.ErrorIs(err, ErrBundleEmptyList)
	})
}

func (r *StorageSQLiteTestSuite) TestSetJournalBundle() {
	r.Run("add food", func() {
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "food_a", Name: "aaa", Brand: "brand a", Cal100: 1, Prot100: 2, Fat100: 3, Carb100: 4, Comment: "Comment",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "food_b", Name: "bbb", Brand: "brand b", Cal100: 5, Prot100: 6, Fat100: 7, Carb100: 8, Comment: "",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "food_c", Name: "ccc", Brand: "brand c", Cal100: 9, Prot100: 10, Fat100: 11, Carb100: 12, Comment: "ccc",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "food_d", Name: "ddd", Brand: "brand d", Cal100: 13, Prot100: 14, Fat100: 15, Carb100: 16, Comment: "ccc",
		}))
		r.NoError(r.stg.SetFood(context.TODO(), &Food{
			Key: "food_e", Name: "eee", Brand: "brand e", Cal100: 17, Prot100: 18, Fat100: 19, Carb100: 20, Comment: "ccc",
		}))
	})

	r.Run("create bundles", func() {
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl1",
			Data: map[string]float64{
				"food_a": 100,
			},
		}))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl2",
			Data: map[string]float64{
				"food_b": 200,
				"bndl1":  0,
			},
		}))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl3",
			Data: map[string]float64{
				"food_c": 300,
				"bndl2":  0,
			},
		}))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl4",
			Data: map[string]float64{
				"food_d": 400,
				"bndl3":  0,
			},
		}))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndl5",
			Data: map[string]float64{
				"food_e": 500,
				"bndl4":  0,
			},
		}))
		//
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndlA",
			Data: map[string]float64{
				"food_a": 100,
				"food_b": 200,
			},
		}))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndlB",
			Data: map[string]float64{
				"food_c": 300,
				"food_d": 400,
			},
		}))
		r.NoError(r.stg.SetBundle(context.TODO(), 1, &Bundle{
			Key: "bndlC",
			Data: map[string]float64{
				"food_e": 500,
				"bndlA":  0,
				"bndlB":  0,
			},
		}))
	})

	r.Run("set journal bundle", func() {
		r.NoError(r.stg.SetJournalBundle(context.TODO(), 1, T(1), Meal(0), "bndl5"))
		r.NoError(r.stg.SetJournalBundle(context.TODO(), 1, T(1), Meal(1), "bndlC"))
	})

	r.Run("check journal", func() {
		rep, err := r.stg.GetJournalReport(context.TODO(), 1, T(1), T(2))
		r.NoError(err)
		r.Equal([]JournalReport{
			{Timestamp: T(1), Meal: Meal(0), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 100, Cal: 1, Prot: 2, Fat: 3, Carb: 4},
			{Timestamp: T(1), Meal: Meal(0), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 200, Cal: 10, Prot: 12, Fat: 14, Carb: 16},
			{Timestamp: T(1), Meal: Meal(0), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 27, Prot: 30, Fat: 33, Carb: 36},
			{Timestamp: T(1), Meal: Meal(0), FoodKey: "food_d", FoodName: "ddd", FoodBrand: "brand d",
				FoodWeight: 400, Cal: 52, Prot: 56, Fat: 60, Carb: 64},
			{Timestamp: T(1), Meal: Meal(0), FoodKey: "food_e", FoodName: "eee", FoodBrand: "brand e",
				FoodWeight: 500, Cal: 85, Prot: 90, Fat: 95, Carb: 100},
			//
			{Timestamp: T(1), Meal: Meal(1), FoodKey: "food_a", FoodName: "aaa", FoodBrand: "brand a",
				FoodWeight: 100, Cal: 1, Prot: 2, Fat: 3, Carb: 4},
			{Timestamp: T(1), Meal: Meal(1), FoodKey: "food_b", FoodName: "bbb", FoodBrand: "brand b",
				FoodWeight: 200, Cal: 10, Prot: 12, Fat: 14, Carb: 16},
			{Timestamp: T(1), Meal: Meal(1), FoodKey: "food_c", FoodName: "ccc", FoodBrand: "brand c",
				FoodWeight: 300, Cal: 27, Prot: 30, Fat: 33, Carb: 36},
			{Timestamp: T(1), Meal: Meal(1), FoodKey: "food_d", FoodName: "ddd", FoodBrand: "brand d",
				FoodWeight: 400, Cal: 52, Prot: 56, Fat: 60, Carb: 64},
			{Timestamp: T(1), Meal: Meal(1), FoodKey: "food_e", FoodName: "eee", FoodBrand: "brand e",
				FoodWeight: 500, Cal: 85, Prot: 90, Fat: 95, Carb: 100},
		}, rep)
	})
}

//
// Activity
//

func (r *StorageSQLiteTestSuite) TestGetActivityList() {
	r.Run("empty list", func() {
		_, err := r.stg.GetActivityList(context.TODO(), 1, T(0), T(10))
		r.ErrorIs(err, ErrActivityEmptyList)
	})

	r.Run("add data", func() {
		r.NoError(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(1), ActiveCal: 1}))
		r.NoError(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(2), ActiveCal: 2}))
		r.NoError(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(3), ActiveCal: 3}))

		r.NoError(r.stg.SetActivity(context.TODO(), 2, &Activity{Timestamp: T(4), ActiveCal: 4}))
	})

	r.Run("get list for different users", func() {
		lst, err := r.stg.GetActivityList(context.TODO(), 1, T(1), T(3))
		r.NoError(err)
		r.Equal([]Activity{
			{Timestamp: T(1), ActiveCal: 1},
			{Timestamp: T(2), ActiveCal: 2},
			{Timestamp: T(3), ActiveCal: 3},
		}, lst)

		lst, err = r.stg.GetActivityList(context.TODO(), 2, T(4), T(4))
		r.NoError(err)
		r.Equal([]Activity{
			{Timestamp: T(4), ActiveCal: 4},
		}, lst)
	})

	r.Run("get limited list", func() {
		lst, err := r.stg.GetActivityList(context.TODO(), 1, T(2), T(3))
		r.NoError(err)
		r.Equal([]Activity{
			{Timestamp: T(2), ActiveCal: 2},
			{Timestamp: T(3), ActiveCal: 3},
		}, lst)
	})
}

func (r *StorageSQLiteTestSuite) TestDeleteActivity() {
	r.Run("add data", func() {
		r.NoError(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(1), ActiveCal: 1}))
		r.NoError(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(2), ActiveCal: 2}))
		r.NoError(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(3), ActiveCal: 3}))

		r.NoError(r.stg.SetActivity(context.TODO(), 2, &Activity{Timestamp: T(4), ActiveCal: 4}))
	})

	r.Run("delete with incorrect user", func() {
		r.NoError(r.stg.DeleteActivity(context.TODO(), 2, T(2)))

		// Data not changed
		lst, err := r.stg.GetActivityList(context.TODO(), 1, T(1), T(3))
		r.NoError(err)
		r.Equal([]Activity{
			{Timestamp: T(1), ActiveCal: 1},
			{Timestamp: T(2), ActiveCal: 2},
			{Timestamp: T(3), ActiveCal: 3},
		}, lst)

		lst, err = r.stg.GetActivityList(context.TODO(), 2, T(4), T(4))
		r.NoError(err)
		r.Equal([]Activity{
			{Timestamp: T(4), ActiveCal: 4},
		}, lst)
	})

	r.Run("delete activity for user", func() {
		r.NoError(r.stg.DeleteActivity(context.TODO(), 2, T(4)))
		_, err := r.stg.GetActivityList(context.TODO(), 2, T(4), T(4))
		r.ErrorIs(err, ErrActivityEmptyList)
	})
}

func (r *StorageSQLiteTestSuite) TestActivityCRU() {
	r.Run("get not existing activity", func() {
		al, err := r.stg.GetActivityList(context.TODO(), 1, T(1), T(5))
		r.ErrorIs(err, ErrActivityEmptyList)
		r.Nil(al)
	})

	r.Run("set invalid activity", func() {
		r.ErrorIs(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(1), ActiveCal: -1}), ErrActivityInvalid)
	})

	r.Run("set valid activity", func() {
		r.NoError(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(1), ActiveCal: 1}))
	})

	r.Run("get activity", func() {
		al, err := r.stg.GetActivityList(context.TODO(), 1, T(1), T(5))
		r.NoError(err)
		r.Equal([]Activity{{Timestamp: T(1), ActiveCal: 1}}, al)

		a, err := r.stg.GetActivity(context.TODO(), 1, T(1))
		r.NoError(err)
		r.Equal(&Activity{Timestamp: T(1), ActiveCal: 1}, a)
	})

	r.Run("set again with update one", func() {
		r.NoError(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(1), ActiveCal: 2}))
		r.NoError(r.stg.SetActivity(context.TODO(), 1, &Activity{Timestamp: T(2), ActiveCal: 2}))
	})

	r.Run("get activity list", func() {
		al, err := r.stg.GetActivityList(context.TODO(), 1, T(1), T(5))
		r.NoError(err)
		r.Equal([]Activity{{Timestamp: T(1), ActiveCal: 2}, {Timestamp: T(2), ActiveCal: 2}}, al)
	})
}

//
// Suite setup
//

func T(sec int) time.Time {
	return time.Date(1970, 1, 1, 0, 0, sec, 0, time.UTC)
}

func (r *StorageSQLiteTestSuite) SetupTest() {
	var err error

	f, err := os.CreateTemp("", "testdb")
	require.NoError(r.T(), err)
	r.dbFile = f.Name()
	f.Close()

	r.stg, err = NewStorageSQLite(r.dbFile)
	// To print queries
	//r.stg, err = NewStorageSQLite(r.dbFile, WithDebug())
	require.NoError(r.T(), err)
}

func (r *StorageSQLiteTestSuite) TearDownTest() {
	r.stg.Close()
	require.NoError(r.T(), os.Remove(r.dbFile))
}

func TestStorageSQLite(t *testing.T) {
	suite.Run(t, new(StorageSQLiteTestSuite))
}
