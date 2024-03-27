package html

import "fmt"

type Script struct {
	url string
}

var _ IELement = (*Script)(nil)

func NewScript(url string) *Script {
	return &Script{url: url}
}

func (r *Script) Build() string {
	return fmt.Sprintf(`<script src="%s"></script>`, r.url)
}
