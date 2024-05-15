package food

import (
	"net/http"

	"github.com/devldavydov/myfood/internal/myfoodserver/templates"
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FoodHandler struct {
	stg    storage.Storage
	logger *zap.Logger
}

func NewFoodHander(stg storage.Storage, logger *zap.Logger) *FoodHandler {
	return &FoodHandler{stg: stg, logger: logger}
}

func (r *FoodHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "food", &templates.TemplateData{
		Nav: "food",
	})
}
