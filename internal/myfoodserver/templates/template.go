package templates

import (
	"embed"
	"html/template"

	"github.com/gin-gonic/gin/render"
)

//go:embed tmpl/*
var fs embed.FS

type TemplateRenderer struct {
	templates map[string]*template.Template
}

func NewTemplateRenderer() *TemplateRenderer {
	templates := map[string]*template.Template{
		"index": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/index.html")),
		"food": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/food.html")),
		"journal": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/journal.html")),
		"weight": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/weight.html")),
		"settings": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/settings.html")),
	}

	return &TemplateRenderer{templates: templates}
}

func (r *TemplateRenderer) Instance(name string, data any) render.Render {
	return render.HTML{
		Template: r.templates[name],
		Data:     data,
	}
}

type TemplateData struct {
	// Error data.
	IsError bool
	Error   string
	// Navigation flag.
	Nav string
	// Custom template data.
	Data any
}
