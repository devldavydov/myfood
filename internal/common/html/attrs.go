package html

import (
	"fmt"
	"strings"
)

type Attrs map[string]string

func (r Attrs) String() string {
	s := make([]string, 0, len(r))
	for k, v := range r {
		s = append(s, fmt.Sprintf(`%s="%s"`, k, v))
	}
	return strings.Join(s, " ")
}
