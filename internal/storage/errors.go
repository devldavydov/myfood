package storage

import "errors"

var (
	// Food
	ErrFoodNotFound      = errors.New("food not found")
	ErrFoodAlreadyExists = errors.New("food already exists")
)
