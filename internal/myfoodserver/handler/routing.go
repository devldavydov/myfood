package handler

import (
	"github.com/devldavydov/myfood/internal/myfoodserver/handler/food"
	"github.com/devldavydov/myfood/internal/myfoodserver/handler/journal"
	"github.com/devldavydov/myfood/internal/myfoodserver/handler/usersettings"
	"github.com/devldavydov/myfood/internal/myfoodserver/handler/weight"
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(rootRouter *gin.Engine, stg storage.Storage, logger *zap.Logger) {
	api := rootRouter.Group("/api")

	food.Attach(api.Group("/food"), stg, logger)
	journal.Attach(api.Group("/journal"), stg, logger)
	usersettings.Attach(api.Group("/usersettings"), stg, logger)
	weight.Attach(api.Group("/weight"), stg, logger)
}
