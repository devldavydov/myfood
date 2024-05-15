package templates

import (
	"embed"
	"fmt"
	"html/template"
	"strings"

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
		"food/list": template.Must(template.ParseFS(fs,
			"tmpl/base.html",
			"tmpl/food/list.html")),
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

type Message struct {
	Cls string
	Msg string
}

func NewMessage(cls, msg string) Message {
	return Message{Cls: cls, Msg: msg}
}

func MessageToString(f Message) string {
	return fmt.Sprintf("%s|%s", f.Cls, f.Msg)
}

func MessageFromString(s string) Message {
	parts := strings.Split(s, "|")
	if len(parts) == 2 {
		return NewMessage(parts[0], parts[1])
	}

	return Message{}
}

type TemplateData struct {
	// Navigation flag.
	Nav string
	// Flash messages.
	Messages []Message
	// Template data.
	Data any
}
