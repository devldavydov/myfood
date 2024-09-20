package storage

import (
	"context"
	"time"
)

type Storage interface {
	// Food
	GetFood(ctx context.Context, key string) (*Food, error)
	SetFood(ctx context.Context, food *Food) error
	SetFoodComment(ctx context.Context, key, comment string) error
	GetFoodList(ctx context.Context) ([]Food, error)
	FindFood(ctx context.Context, pattern string) ([]Food, error)
	DeleteFood(ctx context.Context, key string) error

	// Bundle
	SetBundle(ctx context.Context, userID int64, bndl *Bundle) error
	GetBundle(ctx context.Context, userID int64, key string) (*Bundle, error)
	GetBundleList(ctx context.Context, userID int64) ([]Bundle, error)
	DeleteBundle(ctx context.Context, userID int64, key string) error

	// Weight
	GetWeightList(ctx context.Context, userID int64, from, to time.Time) ([]Weight, error)
	SetWeight(ctx context.Context, userID int64, weight *Weight) error
	DeleteWeight(ctx context.Context, userID int64, timestamp time.Time) error

	// Journal
	SetJournal(ctx context.Context, userID int64, journal *Journal) error
	SetJournalBundle(ctx context.Context, userID int64, timestamp time.Time, meal Meal, bndlKey string) error
	DeleteJournal(ctx context.Context, userID int64, timestamp time.Time, meal Meal, foodkey string) error
	DeleteJournalMeal(ctx context.Context, userID int64, timestamp time.Time, meal Meal) error
	GetJournalMealReport(ctx context.Context, userID int64, timestamp time.Time, meal Meal) (*JournalMealReport, error)
	GetJournalReport(ctx context.Context, userID int64, from, to time.Time) ([]JournalReport, error)
	GetJournalStats(ctx context.Context, userID int64, from, to time.Time) ([]JournalStats, error)
	CopyJournal(ctx context.Context, userID int64, from time.Time, mealFrom Meal, to time.Time, mealTo Meal) (int, error)
	GetJournalFoodAvgWeight(ctx context.Context, userID int64, from, to time.Time, foodkey string) (float64, error)

	// Activity
	GetActivityList(ctx context.Context, userID int64, from, to time.Time) ([]Activity, error)
	GetActivity(ctx context.Context, userID int64, timestamp time.Time) (*Activity, error)
	SetActivity(ctx context.Context, userID int64, activity *Activity) error
	DeleteActivity(ctx context.Context, userID int64, timestamp time.Time) error

	// UserSettings
	GetUserSettings(ctx context.Context, userID int64) (*UserSettings, error)
	SetUserSettings(ctx context.Context, userID int64, settings *UserSettings) error

	Close() error
}
