package storage

import "errors"

var (
	// Food
	ErrFoodNotFound      = errors.New("food not found")
	ErrFoodAlreadyExists = errors.New("food already exists")

	// Weight
	ErrWeightNotFound      = errors.New("weight not found")
	ErrWeightAlreadyExists = errors.New("weight already exists")
	ErrInvalidWeight       = errors.New("invalid weight")
)
