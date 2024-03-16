package storage

import "context"

type Storage interface {
	// Food
	GetFoodList(ctx context.Context, userID int64) ([]Food, error)
	GetFood(ctx context.Context, userID, id int64) (*Food, error)
	FindFood(ctx context.Context, userID int64, query string) ([]Food, error)

	CreateFood(ctx context.Context, userID int64, food *Food) (int64, error)

	DeleteFood(ctx context.Context, userID, id int64) error
	ClearFood(ctx context.Context, userID int64) error
}
