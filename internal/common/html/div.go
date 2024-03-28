package html

import (
	"fmt"
	"strings"
)

type Div struct {
	class    string
	elements []IELement
}

var _ IELement = (*Div)(nil)

func NewDiv(class string) *Div {
	return &Div{class: class}
}

func NewContainer() *Div {
	return NewDiv("container")
}

func (r *Div) Add(elems ...IELement) *Div {
	r.elements = append(r.elements, elems...)
	return r
}

func (r *Div) Build() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`<div class="%s">`, r.class))
	for _, elem := range r.elements {
		sb.WriteString(elem.Build())
	}
	sb.WriteString("</div>")

	return sb.String()
}
