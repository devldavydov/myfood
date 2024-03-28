package html

import "fmt"

type Canvas struct {
	id string
}

var _ IELement = (*Canvas)(nil)

func NewCanvas(id string) *Canvas {
	return &Canvas{id: id}
}

func (r *Canvas) Build() string {
	return fmt.Sprintf(`<canvas id="%s"></canvas>`, r.id)
}
