package weight

import (
	"net/http"

	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WeightHandler struct {
	stg    storage.Storage
	logger *zap.Logger
}

func NewWeightHander(stg storage.Storage, logger *zap.Logger) *WeightHandler {
	return &WeightHandler{stg: stg, logger: logger}
}

func (r *WeightHandler) Dummy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"group": "weight"})
}
