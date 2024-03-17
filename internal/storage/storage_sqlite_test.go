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

func (r *StorageSQLiteTestSuite) TestCreateAndGetWeight() {
	r.Run("create invalid weight", func() {
		w := &Weight{Timestamp: -1, Value: 123}
		r.ErrorIs(r.stg.CreateWeight(context.TODO(), 1, w), ErrWeightInvalid)

		w = &Weight{Timestamp: 0, Value: -1}
		r.ErrorIs(r.stg.CreateWeight(context.TODO(), 1, w), ErrWeightInvalid)
	})

	r.Run("create valid weight", func() {
		w := &Weight{Timestamp: 123, Value: 456}
		r.NoError(r.stg.CreateWeight(context.TODO(), 1, w))
	})

	r.Run("create already exists weight", func() {
		w := &Weight{Timestamp: 123, Value: 123}
		r.ErrorIs(r.stg.CreateWeight(context.TODO(), 1, w), ErrWeightAlreadyExists)
	})

	r.Run("get weight", func() {
		w, err := r.stg.GetWeight(context.TODO(), 1, 123)
		r.NoError(err)
		r.Equal(int64(123), w.Timestamp)
		r.Equal(float64(456), w.Value)
	})

	r.Run("get not existing weight", func() {
		_, err := r.stg.GetWeight(context.TODO(), 111, 111)
		r.ErrorIs(err, ErrWeightNotFound)
	})
}

func (r *StorageSQLiteTestSuite) TestGetWeightList() {
	r.Run("empty list", func() {
		_, err := r.stg.GetWeightList(context.TODO(), 1, 0, 10)
		r.ErrorIs(err, ErrWeightEmptyList)
	})

	r.Run("add data", func() {
		r.NoError(r.stg.CreateWeight(context.TODO(), 1, &Weight{Timestamp: 1, Value: 1}))
		r.NoError(r.stg.CreateWeight(context.TODO(), 1, &Weight{Timestamp: 2, Value: 2}))
		r.NoError(r.stg.CreateWeight(context.TODO(), 1, &Weight{Timestamp: 3, Value: 3}))

		r.NoError(r.stg.CreateWeight(context.TODO(), 2, &Weight{Timestamp: 4, Value: 4}))
	})

	r.Run("get list for different users", func() {
		lst, err := r.stg.GetWeightList(context.TODO(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: 3, Value: 3},
			{Timestamp: 2, Value: 2},
			{Timestamp: 1, Value: 1},
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
			{Timestamp: 3, Value: 3},
			{Timestamp: 2, Value: 2},
		}, lst)
	})
}

func (r *StorageSQLiteTestSuite) TestDeleteAndClearWeight() {
	r.Run("add data", func() {
		r.NoError(r.stg.CreateWeight(context.TODO(), 1, &Weight{Timestamp: 1, Value: 1}))
		r.NoError(r.stg.CreateWeight(context.TODO(), 1, &Weight{Timestamp: 2, Value: 2}))
		r.NoError(r.stg.CreateWeight(context.TODO(), 1, &Weight{Timestamp: 3, Value: 3}))

		r.NoError(r.stg.CreateWeight(context.TODO(), 2, &Weight{Timestamp: 4, Value: 4}))
	})

	r.Run("delete with incorrect user", func() {
		r.NoError(r.stg.DeleteWeight(context.TODO(), 2, 2))

		// Data not changed
		lst, err := r.stg.GetWeightList(context.TODO(), 1, 1, 3)
		r.NoError(err)
		r.Equal([]Weight{
			{Timestamp: 3, Value: 3},
			{Timestamp: 2, Value: 2},
			{Timestamp: 1, Value: 1},
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

	r.Run("clear weight for user", func() {
		r.NoError(r.stg.ClearWeight(context.TODO(), 1))
		_, err := r.stg.GetWeightList(context.TODO(), 1, 0, 10)
		r.ErrorIs(err, ErrWeightEmptyList)
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
