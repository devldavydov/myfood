package food

import (
	"context"
	"errors"
	"net/http"

	"github.com/devldavydov/myfood/internal/common/messages"
	"github.com/devldavydov/myfood/internal/myfoodserver/model"
	"github.com/devldavydov/myfood/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FoodGetResponse struct {
	Key     string  `json:"key"`
	Name    string  `json:"name"`
	Brand   string  `json:"brand"`
	Cal100  float64 `json:"cal100"`
	Prot100 float64 `json:"prot100"`
	Fat100  float64 `json:"fat100"`
	Carb100 float64 `json:"carb100"`
	Comment string  `json:"comment"`
}

func (r *FoodHandler) GetAPI(c *gin.Context) {
	// Get from DB
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, c.Param("key"))
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrFoodNotFound))
			return
		}

		r.logger.Error(
			"food get api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrInternal))
		return
	}

	c.JSON(http.StatusOK, model.NewDataResponse(&FoodGetResponse{
		Key:     food.Key,
		Name:    food.Name,
		Brand:   food.Brand,
		Cal100:  food.Cal100,
		Prot100: food.Prot100,
		Fat100:  food.Fat100,
		Carb100: food.Carb100,
		Comment: food.Comment,
	}))
}

type DeleteFoodRequest struct {
	Key string `json:"key"`
}

func (r *FoodHandler) DeleteAPI(c *gin.Context) {
	req := &DeleteFoodRequest{}
	if err := c.BindJSON(&req); err != nil {
		return
	}

	c.JSON(http.StatusOK, map[string]string{"ok": "true"})
}
