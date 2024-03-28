package html

import "fmt"

type H struct {
	value string
	size  int
	attrs Attrs
}

var _ IELement = (*H)(nil)

func NewH(value string, size int, attrs Attrs) *H {
	return &H{value: value, size: size, attrs: attrs}
}

func (r *H) Build() string {
	return fmt.Sprintf(`<h%d %s>%s</h%d>`, r.size, r.attrs.String(), r.value, r.size)
}
