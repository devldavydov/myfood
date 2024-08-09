package html

import "fmt"

type I struct {
	val   string
	attrs Attrs
}

var _ IELement = (*I)(nil)

func NewI(val string, attrs Attrs) *I {
	return &I{val: val, attrs: attrs}
}

func (r *I) Build() string {
	return fmt.Sprintf("<i %s>%s</i>", r.attrs.String(), r.val)
}
