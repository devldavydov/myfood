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
		// Root
		"index": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/index.html")),
		// Food
		"food/list": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/food/list.html")),
		"food/view": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/food/view.html")),
		// Journal
		"journal": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/journal.html")),
		// Weight
		"weight": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/weight.html")),
		// Settings
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
