package handler

import (
	"github.com/devldavydov/myfood/internal/myfoodserver/handlers/food"
	"github.com/devldavydov/myfood/internal/myfoodserver/handlers/journal"
	"github.com/devldavydov/myfood/internal/myfoodserver/handlers/settings"
	"github.com/devldavydov/myfood/internal/myfoodserver/handlers/weight"
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(router *gin.Engine, stg storage.Storage, logger *zap.Logger) {
	api := router.Group("/api")

	food.Attach(api.Group("/food"), stg, logger)
	journal.Attach(api.Group("/journal"), stg, logger)
	settings.Attach(api.Group("/settings"), stg, logger)
	weight.Attach(api.Group("/weight"), stg, logger)
}
