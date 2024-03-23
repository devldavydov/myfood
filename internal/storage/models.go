package storage

import "strings"

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
	Timestamp int64
	Meal      Meal
	FoodKey   string
	FoodLabel string
	Weight    float64
	Cal       float64
	Prot      float64
	Fat       float64
	Carb      float64
}

func (r *Journal) Validate() bool {
	return r.Timestamp >= 0 &&
		r.Meal >= 0 &&
		r.FoodKey != "" &&
		r.Weight > 0 &&
		r.Cal >= 0 &&
		r.Prot >= 0 &&
		r.Fat >= 0 &&
		r.Carb >= 0
}

type Weight struct {
	Timestamp int64
	Value     float64
}

func (r *Weight) Validate() bool {
	return r.Timestamp >= 0 && r.Value > 0
}

type UserSettings struct {
	CalLimit float64
}

func (r *UserSettings) Validate() bool {
	return r.CalLimit > 0
}
