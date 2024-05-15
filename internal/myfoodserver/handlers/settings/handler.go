package settings

import (
	"net/http"

	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SettingsHandler struct {
	stg    storage.Storage
	logger *zap.Logger
}

func NewSettingsHander(stg storage.Storage, logger *zap.Logger) *SettingsHandler {
	return &SettingsHandler{stg: stg, logger: logger}
}

func (r *SettingsHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "settings", gin.H{"nav": "settings"})
}
