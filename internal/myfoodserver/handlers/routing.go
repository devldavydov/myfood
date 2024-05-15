package handler

import (
	"github.com/devldavydov/myfood/internal/myfoodserver/handlers/food"
	"github.com/devldavydov/myfood/internal/myfoodserver/handlers/journal"
	"github.com/devldavydov/myfood/internal/myfoodserver/handlers/root"
	"github.com/devldavydov/myfood/internal/myfoodserver/handlers/settings"
	"github.com/devldavydov/myfood/internal/myfoodserver/handlers/weight"
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(router *gin.Engine, stg storage.Storage, logger *zap.Logger) {
	root.Attach(router.Group("/"), stg, logger)
	food.Attach(router.Group("/food"), stg, logger)
	journal.Attach(router.Group("/journal"), stg, logger)
	settings.Attach(router.Group("/settings"), stg, logger)
	weight.Attach(router.Group("/weight"), stg, logger)
}
