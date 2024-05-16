package root

import (
	"net/http"

	"github.com/devldavydov/myfood/internal/myfoodserver/model"
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RootHandler struct {
	stg    storage.Storage
	logger *zap.Logger
}

func NewRootHander(stg storage.Storage, logger *zap.Logger) *RootHandler {
	return &RootHandler{stg: stg, logger: logger}
}

func (r *RootHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", &model.TemplateData{
		Nav: "index",
		Data: gin.H{
			"name": "myfood",
		},
	})
}
