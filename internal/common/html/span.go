package html

import "strings"

type Span struct {
	elements []IELement
}

var _ IELement = (*Span)(nil)

func NewSpan(elements ...IELement) *Span {
	return &Span{
		elements: elements,
	}
}

func (r *Span) Build() string {
	var sb strings.Builder
	sb.WriteString("<span>")
	for _, elem := range r.elements {
		sb.WriteString(elem.Build())
	}
	sb.WriteString("</span>")
	return sb.String()
}
