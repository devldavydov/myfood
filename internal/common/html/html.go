package html

import (
	"fmt"
	"strings"
)

const (
	_cssBotstrapURL = "https://devldavydov.github.io/css/bootstrap/bootstrap.min.css"
	_jsBootstrapURL = "https://devldavydov.github.io/js/bootstrap/bootstrap.bundle.min.js"
	_jsChartURL     = "https://devldavydov.github.io/js/chartjs/chart.umd.min.js"
)

type IELement interface {
	Build() string
}

type Builder struct {
	title    string
	elements []IELement
}

func NewBuilder(title string) *Builder {
	return &Builder{title: title}
}

func (r *Builder) Add(elems ...IELement) *Builder {
	r.elements = append(r.elements, elems...)
	return r
}

func (r *Builder) Build() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`
	<!doctype html>
	<html lang="ru">
	
	<head>
	  <meta charset="utf-8">
	  <meta name="viewport" content="width=device-width, initial-scale=1">
	  <title>%s</title>
	  <link href="%s" rel="stylesheet">
	</head>
	<body>
	`,
		r.title,
		_cssBotstrapURL))

	for _, elem := range r.elements {
		sb.WriteString(elem.Build())
	}

	sb.WriteString(`
	</body>
	</html>
	`)

	return sb.String()
}
