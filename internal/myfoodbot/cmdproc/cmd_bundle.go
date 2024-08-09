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

func (r *CmdProcessor) processBundle(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid bundle command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "set":
		resp = r.bundleSetCommand(cmdParts[1:], userID)
	case "st":
		resp = r.bundleSetTemplateCommand(cmdParts[1:], userID)
	case "list":
		resp = r.bundleListCommand(cmdParts[1:], userID)
	case "del":
		resp = r.bundleDelCommand(cmdParts[1:], userID)
	default:
		r.logger.Error(
			"invalid bundle command",
			zap.String("reason", "unknown command"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		resp = NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	return resp
}

func (r *CmdProcessor) bundleSetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) < 1 {
		r.logger.Error(
			"invalid bundle set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	bndlKey := cmdParts[0]
	bndlData := make(map[string]float64)

	for _, cmdPart := range cmdParts[1:] {
		if strings.Contains(cmdPart, ":") {
			// Add dependant food
			parts := strings.Split(cmdPart, ":")
			if len(parts) > 2 {
				return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
			}

			weight, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
			}

			bndlData[parts[0]] = weight
		} else {
			// Add dependant bundle key.
			bndlData[cmdPart] = 0
		}
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetBundle(ctx, userID, &storage.Bundle{Key: bndlKey, Data: bndlData}); err != nil {
		if errors.Is(err, storage.ErrBundleInvalid) {
			return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
		}
		if errors.Is(err, storage.ErrBundleDepBundleNotFound) {
			return NewSingleCmdResponse(messages.MsgErrBundleDepBundleNotFound)
		}
		if errors.Is(err, storage.ErrBundleDepRecursive) {
			return NewSingleCmdResponse(messages.MsgErrBundleDepBundleRecursive)
		}
		if errors.Is(err, storage.ErrBundleDepFoodNotFound) {
			return NewSingleCmdResponse(messages.MsgErrBundleDepFoodNotFound)
		}

		r.logger.Error(
			"bundle set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) bundleSetTemplateCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid bundle set template command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	bndl, err := r.stg.GetBundle(ctx, userID, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrBundleNotFound) {
			return NewSingleCmdResponse(messages.MsgErrBundleNotFound)
		}

		r.logger.Error(
			"bundle set template command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("b,set,%s", bndl.Key))
	for k, v := range bndl.Data {
		sb.WriteString(",")
		if v > 0 {
			sb.WriteString(fmt.Sprintf("%s:%1.f", k, v))
		} else {
			sb.WriteString(k)
		}
	}

	return NewSingleCmdResponse(sb.String())
}

func (r *CmdProcessor) bundleListCommand(cmdParts []string, userID int64) []CmdResponse {
	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	lst, err := r.stg.GetBundleList(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrBundleEmptyList) {
			return NewSingleCmdResponse(messages.MsgErrEmptyList)
		}

		r.logger.Error(
			"bundle list command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	// Build html
	htmlBuilder := html.NewBuilder("Список бандлов")

	// Table
	tbl := html.NewTable([]string{
		"Ключ бандла", "Еда/Ключ дочернего бандла", "Вес еды, г.",
	})

	for _, bndl := range lst {
		i := 0
		for k, v := range bndl.Data {
			tr := html.NewTr(nil)
			if i == 0 {
				tr.AddTd(html.NewTd(html.NewS(bndl.Key), html.Attrs{"rowspan": strconv.Itoa(len(bndl.Data))}))
			}
			if v > 0 {
				tr.AddTd(html.NewTd(html.NewS(k), nil))
				tr.AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.1f", v)), nil))
			} else {
				tr.AddTd(html.NewTd(html.NewI(k, nil), nil))
				tr.AddTd(html.NewTd(html.NewS(""), nil))
			}
			i++
			tbl.AddRow(tr)
		}
	}

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			html.NewH(
				"Список бандлов",
				5,
				html.Attrs{"align": "center"},
			),
			tbl))

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: "bundles.html",
	})
}

func (r *CmdProcessor) bundleDelCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid bundle del command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteBundle(ctx, userID, cmdParts[0]); err != nil {
		if errors.Is(err, storage.ErrBundleIsUsed) {
			return NewSingleCmdResponse(messages.MsgErrBundleIsUsed)
		}

		r.logger.Error(
			"bundle del command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}
