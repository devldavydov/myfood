package storage

import "context"

type Storage interface {
	// // Food
	// GetFoodList(ctx context.Context, userID int64) ([]Food, error)
	// GetFood(ctx context.Context, userID, id int64) (*Food, error)
	// FindFood(ctx context.Context, userID int64, query string) ([]Food, error)

	// CreateFood(ctx context.Context, userID int64, food *Food) (int64, error)
	// UpdateFood(ctx context.Context, userID int64, food *Food) error

	// DeleteFood(ctx context.Context, userID, id int64) error
	// ClearFood(ctx context.Context, userID int64) error

	// Weight
	GetWeightList(ctx context.Context, userID int64, from, to int64) ([]Weight, error)
	GetWeight(ctx context.Context, userID int64, timestamp int64) (*Weight, error)

	CreateWeight(ctx context.Context, userID int64, weight *Weight) error
	UpdateWeight(ctx context.Context, userID int64, weight *Weight) error

	DeleteWeight(ctx context.Context, userID, timestamp int64) error
	ClearWeight(ctx context.Context, userID int64) error

	Close()
}
