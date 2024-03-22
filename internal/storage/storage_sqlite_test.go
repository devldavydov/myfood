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
