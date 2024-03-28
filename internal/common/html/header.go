package html

import "fmt"

type Header struct {
	value string
	size  int
	attrs Attrs
}

var _ IELement = (*Header)(nil)

func NewHeader(value string, size int, attrs Attrs) *Header {
	return &Header{value: value, size: size, attrs: attrs}
}

func (r *Header) Build() string {
	return fmt.Sprintf(`<h%d %s>%s/h%d>`, r.size, r.attrs.String(), r.value, r.size)
}
