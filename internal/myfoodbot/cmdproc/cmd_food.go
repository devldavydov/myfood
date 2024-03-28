package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/devldavydov/myfood/internal/common/html"
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) processFood(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid food command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	switch cmdParts[0] {
	case "set":
		return r.foodSetCommand(c, cmdParts[1:], userID)
	case "sc":
		return r.foodSetCommentCommand(c, cmdParts[1:], userID)
	case "find":
		return r.foodFindCommand(c, cmdParts[1:], userID)
	case "list":
		return r.foodListCommand(c, userID)
	case "del":
		return r.foodDelCommand(c, cmdParts[1:], userID)
	}

	r.logger.Error(
		"invalid food command",
		zap.String("reason", "unknown command"),
		zap.Strings("command", cmdParts),
		zap.Int64("userid", userID),
	)
	return c.Send(msgErrInvalidCommand)
}

func (r *CmdProcessor) foodSetCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 8 {
		r.logger.Error(
			"invalid food set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
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
		return c.Send(msgErrInvalidCommand)
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
		return c.Send(msgErrInvalidCommand)
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
		return c.Send(msgErrInvalidCommand)
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
		return c.Send(msgErrInvalidCommand)
	}
	food.Carb100 = carb100

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.SetFood(ctx, food); err != nil {
		if errors.Is(err, storage.ErrFoodInvalid) {
			return c.Send(msgErrInvalidCommand)
		}

		r.logger.Error(
			"food set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *CmdProcessor) foodSetCommentCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid food set comment command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.SetFoodComment(ctx, cmdParts[0], cmdParts[1]); err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return c.Send(msgErrFoodNotFound)
		}

		r.logger.Error(
			"food set comment command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *CmdProcessor) foodFindCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid food find command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	foodLst, err := r.stg.FindFood(ctx, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrFoodEmptyList) {
			return c.Send(msgErrEmptyList)
		}

		r.logger.Error(
			"food find command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
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

	return c.Send(sb.String(), &tele.SendOptions{ParseMode: tele.ModeHTML})
}

func (r *CmdProcessor) foodDelCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid food del command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteFood(ctx, cmdParts[0]); err != nil {
		if errors.Is(err, storage.ErrFoodIsUsed) {
			return c.Send(msgErrFoodIsUsed)
		}

		r.logger.Error(
			"food del command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *CmdProcessor) foodListCommand(c tele.Context, userID int64) error {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	foodList, err := r.stg.GetFoodList(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrFoodEmptyList) {
			return c.Send(msgErrEmptyList)
		}

		r.logger.Error(
			"food list command DB error",
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
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
			AddTd(html.NewTd(html.NewS(item.Cal100), nil)).
			AddTd(html.NewTd(html.NewS(item.Prot100), nil)).
			AddTd(html.NewTd(html.NewS(item.Fat100), nil)).
			AddTd(html.NewTd(html.NewS(item.Carb100), nil)).
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
	return c.Send(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: "food.html",
	})
}
