package helpers

import (
	"fmt"

	"github.com/devldavydov/myfood/internal/myfoodserver/templates"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AddFlashMessage(c *gin.Context, cls, msg string) {
	sess := sessions.Default(c)
	sess.AddFlash(fmt.Sprintf("%s|%s", cls, msg))
	sess.Save()
}

func InitTemplateData(c *gin.Context, nav string) *templates.TemplateData {
	sess := sessions.Default(c)
	var msgs []templates.Message
	for _, f := range sess.Flashes() {
		msgs = append(msgs, templates.MessageFromString(f.(string)))
	}
	sess.Save()

	return &templates.TemplateData{
		Nav:      nav,
		Messages: msgs,
	}
}
