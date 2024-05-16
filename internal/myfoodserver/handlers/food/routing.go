package food

import (
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Attach(group *gin.RouterGroup, stg storage.Storage, logger *zap.Logger) {
	foodHandler := NewFoodHander(stg, logger)

	// Pages
	group.GET("/", foodHandler.IndexPage)
	group.GET("/view/:key", foodHandler.ViewPage)
	group.GET("/set/:key", foodHandler.SetPage)

	// API
	api := group.Group("/api")
	api.GET("/get/:key", foodHandler.GetAPI)
	api.POST("/del", foodHandler.DeleteAPI)
}
