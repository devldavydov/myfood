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

type FoodHandler struct {
	stg    storage.Storage
	logger *zap.Logger
}

func NewFoodHander(stg storage.Storage, logger *zap.Logger) *FoodHandler {
	return &FoodHandler{stg: stg, logger: logger}
}

type FoodResponseItem struct {
	Key     string  `json:"key"`
	Name    string  `json:"name"`
	Brand   string  `json:"brand"`
	Cal100  float64 `json:"cal100"`
	Prot100 float64 `json:"prot100"`
	Fat100  float64 `json:"fat100"`
	Carb100 float64 `json:"carb100"`
	Comment string  `json:"comment"`
}

func (r *FoodHandler) ListAPI(c *gin.Context) {
	// Get from DB
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	foodList, err := r.stg.GetFoodList(ctx)
	if err != nil && !errors.Is(err, storage.ErrFoodEmptyList) {
		r.logger.Error(
			"food list api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrInternal))
		return
	}

	data := make([]FoodResponseItem, 0, len(foodList))
	for _, f := range foodList {
		data = append(data, FoodResponseItem{
			Key:     f.Key,
			Name:    f.Name,
			Brand:   f.Brand,
			Cal100:  f.Cal100,
			Comment: f.Comment,
		})
	}

	c.JSON(http.StatusOK, model.NewDataResponse(data))
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

	c.JSON(http.StatusOK, model.NewDataResponse(&FoodResponseItem{
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

	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteFood(ctx, req.Key); err != nil {
		if errors.Is(err, storage.ErrFoodIsUsed) {
			c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrFoodIsUsed))
			return
		}

		r.logger.Error(
			"food del api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrInternal))
	}

	c.JSON(http.StatusOK, model.NewOKResponse())
}
