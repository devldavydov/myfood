package storage

import "errors"

var (
	// Food
	ErrFoodNotFound  = errors.New("food not found")
	ErrFoodInvalid   = errors.New("invalid food")
	ErrFoodEmptyList = errors.New("empty food list")

	// Journal
	ErrJournalInvalid     = errors.New("journal invalid")
	ErrJournalReportEmpty = errors.New("empty journal report")
	ErrJournalInvalidFood = errors.New("journal invalid food")

	// Weight
	ErrWeightNotFound  = errors.New("weight not found")
	ErrWeightInvalid   = errors.New("invalid weight")
	ErrWeightEmptyList = errors.New("empty weight list")

	// UserSettings
	ErrUserSettingsNotFound = errors.New("user settings not found")
	ErrUserSettingsInvalid  = errors.New("invalid user settings")
)
