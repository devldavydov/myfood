package usersettings

import (
	"net/http"

	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserSettingsHandler struct {
	stg    storage.Storage
	logger *zap.Logger
}

func NewUserSettingsHander(stg storage.Storage, logger *zap.Logger) *UserSettingsHandler {
	return &UserSettingsHandler{stg: stg, logger: logger}
}

func (r *UserSettingsHandler) Dummy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"group": "user settings"})
}
