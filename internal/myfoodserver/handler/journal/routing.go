package journal

import (
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Attach(group *gin.RouterGroup, stg storage.Storage, logger *zap.Logger) {
	journalHandler := NewJournalHander(stg, logger)
	group.GET("/", journalHandler.Dummy)
}
