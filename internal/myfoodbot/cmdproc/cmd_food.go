package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/devldavydov/myfood/internal/common/html"
	"github.com/devldavydov/myfood/internal/common/messages"
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) processFood(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid food command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "set":
		resp = r.foodSetCommand(cmdParts[1:], userID)
	case "sc":
		resp = r.foodSetCommentCommand(cmdParts[1:], userID)
	case "st":
		resp = r.foodSetTemplateCommand(cmdParts[1:], userID)
	case "find":
		resp = r.foodFindCommand(cmdParts[1:], userID)
	case "calc":
		resp = r.foodCalcCommand(cmdParts[1:], userID)
	case "list":
		resp = r.foodListCommand(userID)
	case "del":
		resp = r.foodDelCommand(cmdParts[1:], userID)
	default:
		r.logger.Error(
			"invalid food command",
			zap.String("reason", "unknown command"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		resp = NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	return resp
}

func (r *CmdProcessor) foodSetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 8 {
		r.logger.Error(
			"invalid food set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Parse fields
	food := &storage.Food{
		Key:     cmdParts[0],
		Name:    cmdParts[1],
		Brand:   cmdParts[2],
		Comment: cmdParts[7],
	}

	cal100, err := strconv.ParseFloat(cmdParts[3], 64)
	if err != nil {
		r.logger.Error(
			"invalid food set command",
			zap.String("reason", "cal100 format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}
	food.Cal100 = cal100

	prot100, err := strconv.ParseFloat(cmdParts[4], 64)
	if err != nil {
		r.logger.Error(
			"invalid food set command",
			zap.String("reason", "prot100 format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}
	food.Prot100 = prot100

	fat100, err := strconv.ParseFloat(cmdParts[5], 64)
	if err != nil {
		r.logger.Error(
			"invalid food set command",
			zap.String("reason", "fat100 format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}
	food.Fat100 = fat100

	carb100, err := strconv.ParseFloat(cmdParts[6], 64)
	if err != nil {
		r.logger.Error(
			"invalid food set command",
			zap.String("reason", "carb100 format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}
	food.Carb100 = carb100

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetFood(ctx, food); err != nil {
		if errors.Is(err, storage.ErrFoodInvalid) {
			return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
		}

		r.logger.Error(
			"food set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) foodSetCommentCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid food set comment command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetFoodComment(ctx, cmdParts[0], cmdParts[1]); err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(messages.MsgErrFoodNotFound)
		}

		r.logger.Error(
			"food set comment command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) foodSetTemplateCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid food set template command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Get food from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(messages.MsgErrFoodNotFound)
		}

		r.logger.Error(
			"food set template command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	foodSetTemplate := fmt.Sprintf(
		"f,set,%s,%s,%s,%.2f,%.2f,%.2f,%.2f,%s",
		food.Key,
		food.Name,
		food.Brand,
		food.Cal100,
		food.Prot100,
		food.Fat100,
		food.Carb100,
		food.Comment,
	)
	return NewSingleCmdResponse(foodSetTemplate, optsHTML)
}

func (r *CmdProcessor) foodFindCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid food find command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	foodLst, err := r.stg.FindFood(ctx, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrFoodEmptyList) {
			return NewSingleCmdResponse(messages.MsgErrEmptyList)
		}

		r.logger.Error(
			"food find command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	var sb strings.Builder

	for i, food := range foodLst {
		sb.WriteString(fmt.Sprintf("<b>Ключ:</b> %s\n", food.Key))
		sb.WriteString(fmt.Sprintf("<b>Наименование:</b> %s\n", food.Name))
		sb.WriteString(fmt.Sprintf("<b>Бренд:</b> %s\n", food.Brand))
		sb.WriteString(fmt.Sprintf("<b>ККал100:</b> %.2f\n", food.Cal100))
		sb.WriteString(fmt.Sprintf("<b>Бел100:</b> %.2f\n", food.Prot100))
		sb.WriteString(fmt.Sprintf("<b>Жир100:</b> %.2f\n", food.Fat100))
		sb.WriteString(fmt.Sprintf("<b>Угл100:</b> %.2f\n", food.Carb100))
		sb.WriteString(fmt.Sprintf("<b>Комментарий:</b> %s\n", food.Comment))

		if i != len(foodLst)-1 {
			sb.WriteString("\n")
		}
	}

	return NewSingleCmdResponse(sb.String(), optsHTML)
}

func (r *CmdProcessor) foodCalcCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid food calc command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	foodWeight, err := strconv.ParseFloat(cmdParts[1], 64)
	if err != nil {
		r.logger.Error(
			"invalid food calc command",
			zap.String("reason", "weight format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(messages.MsgErrFoodNotFound)
		}

		r.logger.Error(
			"food calc command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>Наименование:</b> %s\n", food.Name))
	sb.WriteString(fmt.Sprintf("<b>Бренд:</b> %s\n", food.Brand))
	sb.WriteString(fmt.Sprintf("<b>Вес:</b> %.1f\n", foodWeight))
	sb.WriteString(fmt.Sprintf("<b>ККал:</b> %.2f\n", foodWeight/100*food.Cal100))
	sb.WriteString(fmt.Sprintf("<b>Бел:</b> %.2f\n", foodWeight/100*food.Prot100))
	sb.WriteString(fmt.Sprintf("<b>Жир:</b> %.2f\n", foodWeight/100*food.Fat100))
	sb.WriteString(fmt.Sprintf("<b>Угл:</b> %.2f\n", foodWeight/100*food.Carb100))

	return NewSingleCmdResponse(sb.String(), optsHTML)
}

func (r *CmdProcessor) foodDelCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid food del command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteFood(ctx, cmdParts[0]); err != nil {
		if errors.Is(err, storage.ErrFoodIsUsed) {
			return NewSingleCmdResponse(messages.MsgErrFoodIsUsed)
		}

		r.logger.Error(
			"food del command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) foodListCommand(userID int64) []CmdResponse {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	foodList, err := r.stg.GetFoodList(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrFoodEmptyList) {
			return NewSingleCmdResponse(messages.MsgErrEmptyList)
		}

		r.logger.Error(
			"food list command DB error",
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	// Build html
	htmlBuilder := html.NewBuilder("Список продуктов")

	// Table
	tbl := html.NewTable([]string{
		"Ключ", "Наименование", "Бренд", "ККал в 100г.", "Белки в 100г.",
		"Жиры в 100г.", "Углеводы в 100г.", "Комментарий",
	})

	for _, item := range foodList {
		tr := html.NewTr(nil)
		tr.
			AddTd(html.NewTd(html.NewS(item.Key), nil)).
			AddTd(html.NewTd(html.NewS(item.Name), nil)).
			AddTd(html.NewTd(html.NewS(item.Brand), nil)).
			AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", item.Cal100)), nil)).
			AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", item.Prot100)), nil)).
			AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", item.Fat100)), nil)).
			AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", item.Carb100)), nil)).
			AddTd(html.NewTd(html.NewS(item.Comment), nil))
		tbl.AddRow(tr)
	}

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			html.NewH(
				"Список продуктов и энергетической ценности",
				5,
				html.Attrs{"align": "center"},
			),
			tbl))

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: "food.html",
	})
}
