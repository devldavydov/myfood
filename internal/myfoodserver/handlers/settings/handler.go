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
	c.JSON(http.StatusOK, gin.H{"settings": true})
}
