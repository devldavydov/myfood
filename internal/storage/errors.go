package storage

import "errors"

var (
	// Food
	ErrFoodNotFound  = errors.New("food not found")
	ErrFoodInvalid   = errors.New("invalid food")
	ErrFoodEmptyList = errors.New("empty food list")
	ErrFoodIsUsed    = errors.New("food is used")

	// Bundle
	ErrBundleInvalid           = errors.New("invalid bundle")
	ErrBundleNotFound          = errors.New("bundle not found")
	ErrBundleEmptyList         = errors.New("empty bundle list")
	ErrBundleDepBundleNotFound = errors.New("dependent bundle not found")
	ErrBundleDepFoodNotFound   = errors.New("dependent food not found")
	ErrBundleDepRecursive      = errors.New("dependent recursive bundle not allowed")
	ErrBundleIsUsed            = errors.New("bundle is used")

	// Journal
	ErrJournalInvalid         = errors.New("journal invalid")
	ErrJournalMealReportEmpty = errors.New("empty journal meal report")
	ErrJournalReportEmpty     = errors.New("empty journal report")
	ErrJournalStatsEmpty      = errors.New("empty journal stats")
	ErrJournalInvalidFood     = errors.New("journal invalid food")
	ErrCopyToNotEmpty         = errors.New("copy destination not empty")

	// Weight
	ErrWeightNotFound  = errors.New("weight not found")
	ErrWeightInvalid   = errors.New("invalid weight")
	ErrWeightEmptyList = errors.New("empty weight list")

	// UserSettings
	ErrUserSettingsNotFound = errors.New("user settings not found")
	ErrUserSettingsInvalid  = errors.New("invalid user settings")
)
