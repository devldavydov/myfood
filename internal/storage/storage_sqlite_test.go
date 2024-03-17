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
		r.ErrorIs(r.stg.CreateWeight(context.TODO(), 1, w), ErrInvalidWeight)

		w = &Weight{Timestamp: 0, Value: -1}
		r.ErrorIs(r.stg.CreateWeight(context.TODO(), 1, w), ErrInvalidWeight)
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
