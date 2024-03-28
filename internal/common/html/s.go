package html

import "strconv"

type S struct {
	val string
}

var _ IELement = (*S)(nil)

func NewS(val any) *S {
	switch v := val.(type) {
	case string:
		return &S{val: v}
	case int64:
		return &S{val: strconv.FormatInt(v, 10)}
	case float64:
		return &S{val: strconv.FormatFloat(v, 'f', -1, 64)}
	}

	return &S{val: ""}
}

func (r *S) Build() string {
	return r.val
}
