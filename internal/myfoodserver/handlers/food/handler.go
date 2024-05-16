package food

import (
	"context"
	"errors"
	"net/http"

	"github.com/devldavydov/myfood/internal/common/messages"
	"github.com/devldavydov/myfood/internal/myfoodserver/helpers"
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
	tmplData := helpers.InitTemplateData(c, "food")

	// Get from DB
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	foodList, err := r.stg.GetFoodList(ctx)
	if err != nil && !errors.Is(err, storage.ErrFoodEmptyList) {
		r.logger.Error(
			"food list DB error",
			zap.Error(err),
		)

		tmplData.Messages = append(tmplData.Messages, templates.NewMessage(messages.MsgClassError, messages.MsgErrInternal))
		c.HTML(http.StatusOK, "food/list", tmplData)
		return
	}

	data := make([]map[string]any, 0, len(foodList))
	for _, f := range foodList {
		data = append(data, map[string]any{
			"key":     f.Key,
			"name":    f.Name,
			"brand":   f.Brand,
			"cal100":  f.Cal100,
			"prot100": f.Prot100,
			"fat100":  f.Fat100,
			"carb100": f.Carb100,
			"comment": f.Comment,
		})
	}
	tmplData.Data = data

	c.HTML(http.StatusOK, "food/list", tmplData)
}

func (r *FoodHandler) View(c *gin.Context) {
	tmplData := helpers.InitTemplateData(c, "food")

	// Get from DB
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, c.Param("key"))
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			c.HTML(http.StatusOK, "food/view", tmplData)
			return
		}

		r.logger.Error(
			"food view DB error",
			zap.Error(err),
		)

		tmplData.Messages = append(tmplData.Messages, templates.NewMessage(messages.MsgClassError, messages.MsgErrInternal))
		c.HTML(http.StatusOK, "food/view", tmplData)
		return
	}

	tmplData.Data = map[string]any{
		"key":     food.Key,
		"name":    food.Name,
		"brand":   food.Brand,
		"cal100":  food.Cal100,
		"prot100": food.Prot100,
		"fat100":  food.Fat100,
		"carb100": food.Carb100,
		"comment": food.Comment,
	}

	c.HTML(http.StatusOK, "food/view", tmplData)
}

func (r *FoodHandler) Set(c *gin.Context) {
	helpers.AddFlashMessage(c, messages.MsgClassWarning, messages.MsgErrUnderCon)
	helpers.AddFlashMessage(c, messages.MsgClassWarning, messages.MsgErrUnderCon)
	helpers.AddFlashMessage(c, messages.MsgClassWarning, messages.MsgErrUnderCon)
	c.Redirect(http.StatusMovedPermanently, "/food")
}
