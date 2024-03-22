package storage

type Food struct {
	ID      int64
	Name    string
	Brand   string
	Cal100  float64
	Prot100 float64
	Fat100  float64
	Carb100 float64
	Comment string
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
