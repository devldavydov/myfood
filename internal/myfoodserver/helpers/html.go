package helpers

import "github.com/gin-gonic/gin"

func SendHTML(c *gin.Context, code int, data []byte) {
	c.Data(code, "text/html; charset=utf-8", data)
}
