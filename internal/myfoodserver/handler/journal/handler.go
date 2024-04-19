package journal

import (
	"net/http"

	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type JournalHandler struct {
	stg    storage.Storage
	logger *zap.Logger
}

func NewJournalHander(stg storage.Storage, logger *zap.Logger) *JournalHandler {
	return &JournalHandler{stg: stg, logger: logger}
}

func (r *JournalHandler) Dummy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"group": "journal"})
}
