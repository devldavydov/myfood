package storage

import (
	"strings"
	"time"
)

type Food struct {
	Key     string
	Name    string
	Brand   string
	Cal100  float64
	Prot100 float64
	Fat100  float64
	Carb100 float64
	Comment string
}

func (r *Food) Validate() bool {
	return r.Key != "" &&
		r.Name != "" &&
		r.Cal100 >= 0 &&
		r.Prot100 >= 0 &&
		r.Fat100 >= 0 &&
		r.Carb100 >= 0
}

type Meal int64

func NewMealFromString(m string) Meal {
	switch strings.ToUpper(m) {
	case "ЗАВТРАК":
		return Meal(0)
	case "ДО ОБЕДА":
		return Meal(1)
	case "ОБЕД":
		return Meal(2)
	case "ПОЛДНИК":
		return Meal(3)
	case "ДО УЖИНА":
		return Meal(4)
	case "УЖИН":
		return Meal(5)
	}
	return Meal(6)
}

func (r Meal) ToString() string {
	switch r {
	case 0:
		return "Завтрак"
	case 1:
		return "До обеда"
	case 2:
		return "Обед"
	case 3:
		return "Полдник"
	case 4:
		return "До ужина"
	case 5:
		return "Ужин"
	}
	return "Перекус"
}

type Journal struct {
	Timestamp  time.Time
	Meal       Meal
	FoodKey    string
	FoodWeight float64
}

func (r *Journal) Validate() bool {
	return r.Meal >= 0 &&
		r.FoodKey != "" &&
		r.FoodWeight > 0
}

type JournalMealReport struct {
	ConsumedDayCal  float64
	ConsumedMealCal float64
	Items           []JournalMealItem
}

type JournalMealItem struct {
	Timestamp  time.Time
	FoodKey    string
	FoodName   string
	FoodBrand  string
	FoodWeight float64
	Cal        float64
}

type JournalReport struct {
	Timestamp  time.Time
	Meal       Meal
	FoodKey    string
	FoodName   string
	FoodBrand  string
	FoodWeight float64
	Cal        float64
	Prot       float64
	Fat        float64
	Carb       float64
}

type JournalStats struct {
	Timestamp time.Time
	TotalCal  float64
	TotalProt float64
	TotalFat  float64
	TotalCarb float64
}

type Weight struct {
	Timestamp time.Time
	Value     float64
}

func (r *Weight) Validate() bool {
	return r.Value > 0
}

type UserSettings struct {
	CalLimit         float64
	DefaultActiveCal float64
}

func (r *UserSettings) Validate() bool {
	return r.CalLimit > 0 && r.DefaultActiveCal > 0
}

type Bundle struct {
	Key string
	// Map of bundle data
	// Variants:
	// if food: food_key -> weight > 0
	// if bundle: bundle_key -> 0
	Data map[string]float64
}

func (r *Bundle) Validate() bool {
	if r.Key == "" || len(r.Data) == 0 {
		return false
	}

	for _, v := range r.Data {
		if v < 0 {
			return false
		}
	}

	return true
}

type Activity struct {
	Timestamp time.Time
	ActiveCal float64
}

func (r *Activity) Validate() bool {
	return r.ActiveCal > 0
}

type Backup struct {
	Timestamp    int64                `json:"timestamp"`
	Weight       []WeightBackup       `json:"weight"`
	Food         []FoodBackup         `json:"food"`
	Journal      []JournalBackup      `json:"journal"`
	Bundle       []BundleBackup       `json:"bundle"`
	UserSettings []UserSettingsBackup `json:"user_settings"`
}

type WeightBackup struct {
	UserID    int64   `json:"user_id"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type FoodBackup struct {
	Key     string  `json:"key"`
	Name    string  `json:"name"`
	Brand   string  `json:"brand"`
	Cal100  float64 `json:"cal100"`
	Prot100 float64 `json:"prot100"`
	Fat100  float64 `json:"fat100"`
	Carb100 float64 `json:"carb100"`
	Comment string  `json:"comment"`
}

type JournalBackup struct {
	UserID     int64   `json:"user_id"`
	Timestamp  int64   `json:"timestamp"`
	Meal       int64   `json:"meal"`
	FoodKey    string  `json:"food_key"`
	FoodWeight float64 `json:"food_weight"`
}

type BundleBackup struct {
	UserID int64              `json:"user_id"`
	Key    string             `json:"key"`
	Data   map[string]float64 `json:"data"`
}

type UserSettingsBackup struct {
	UserID   int64   `json:"user_id"`
	CalLimit float64 `json:"cal_limit"`
}
