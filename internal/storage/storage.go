package storage

import "context"

type Storage interface {
	// // Food
	GetFood(ctx context.Context, key string) (*Food, error)
	SetFood(ctx context.Context, food *Food) error
	GetFoodList(ctx context.Context) ([]Food, error)
	FindFood(ctx context.Context, pattern string) ([]Food, error)
	DeleteFood(ctx context.Context, key string) error

	// Weight
	GetWeightList(ctx context.Context, userID int64, from, to int64) ([]Weight, error)
	SetWeight(ctx context.Context, userID int64, weight *Weight) error
	DeleteWeight(ctx context.Context, userID, timestamp int64) error

	// Journal
	SetJournal(ctx context.Context, userID int64, journal *Journal) error
	DeleteJournal(ctx context.Context, userID int64, timestamp int64, meal Meal, foodkey string) error
	GetJournalForPeriod(ctx context.Context, userID int64, from, to int64) ([]JournalReport, error)

	// UserSettings
	GetUserSettings(ctx context.Context, userID int64) (*UserSettings, error)
	SetUserSettings(ctx context.Context, userID int64, settings *UserSettings) error

	Close()
}
