package storage

import (
	"context"
	"time"
)

type Storage interface {
	// // Food
	GetFood(ctx context.Context, key string) (*Food, error)
	SetFood(ctx context.Context, food *Food) error
	SetFoodComment(ctx context.Context, key, comment string) error
	GetFoodList(ctx context.Context) ([]Food, error)
	FindFood(ctx context.Context, pattern string) ([]Food, error)
	DeleteFood(ctx context.Context, key string) error

	// Weight
	GetWeightList(ctx context.Context, userID int64, from, to time.Time) ([]Weight, error)
	SetWeight(ctx context.Context, userID int64, weight *Weight) error
	DeleteWeight(ctx context.Context, userID int64, timestamp time.Time) error

	// Journal
	SetJournal(ctx context.Context, userID int64, journal *Journal) error
	DeleteJournal(ctx context.Context, userID int64, timestamp time.Time, meal Meal, foodkey string) error
	GetJournalReport(ctx context.Context, userID int64, from, to time.Time) ([]JournalReport, error)
	GetJournalStats(ctx context.Context, userID int64, from, to time.Time) ([]JournalStats, error)

	// UserSettings
	GetUserSettings(ctx context.Context, userID int64) (*UserSettings, error)
	SetUserSettings(ctx context.Context, userID int64, settings *UserSettings) error

	Close() error
}
