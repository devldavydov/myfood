package html

import "fmt"

type B struct {
	val   string
	attrs Attrs
}

var _ IELement = (*B)(nil)

func NewB(val string, attrs Attrs) *B {
	return &B{val: val, attrs: attrs}
}

func (r *B) Build() string {
	return fmt.Sprintf("<b %s>%s</b>", r.attrs.String(), r.val)
}
