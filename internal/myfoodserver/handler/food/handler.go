package food

import (
	"net/http"

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

func (r *FoodHandler) Dummy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"group": "food"})
}
