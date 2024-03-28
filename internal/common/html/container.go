package html

import "strings"

type Container struct {
	elements []IELement
}

var _ IELement = (*Container)(nil)

func NewContainer() *Container {
	return &Container{}
}

func (r *Container) Add(elems ...IELement) *Container {
	r.elements = append(r.elements, elems...)
	return r
}

func (r *Container) Build() string {
	var sb strings.Builder
	for _, elem := range r.elements {
		sb.WriteString(elem.Build())
	}
	return sb.String()
}
