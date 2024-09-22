package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/devldavydov/myfood/internal/common/html"
	"github.com/devldavydov/myfood/internal/common/messages"
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) processActivity(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid activity command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "set":
		resp = r.activitySetCommand(cmdParts[1:], userID)
	case "list":
		resp = r.activityListCommand(cmdParts[1:], userID)
	case "del":
		resp = r.activityDelCommand(cmdParts[1:], userID)
	default:
		r.logger.Error(
			"invalid activity command",
			zap.String("reason", "unknown command"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		resp = NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	return resp
}

func (r *CmdProcessor) activitySetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid activity set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Parse timestamp
	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid activity set command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Parse activeCal
	activeCal, err := strconv.ParseFloat(cmdParts[1], 64)
	if err != nil {
		r.logger.Error(
			"invalid activity set command",
			zap.String("reason", "activeCal format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetActivity(ctx, userID, &storage.Activity{Timestamp: ts, ActiveCal: activeCal}); err != nil {
		if errors.Is(err, storage.ErrActivityInvalid) {
			return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
		}

		r.logger.Error(
			"activity set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) activityListCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid activity list command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Parse timestamp
	tsFrom, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid activity list command",
			zap.String("reason", "ts from format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	tsTo, err := r.parseTimestamp(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid activity list command",
			zap.String("reason", "ts to format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// List from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	lst, err := r.stg.GetActivityList(ctx, userID, tsFrom, tsTo)
	if err != nil {
		if errors.Is(err, storage.ErrActivityEmptyList) {
			return NewSingleCmdResponse(messages.MsgErrEmptyList)
		}

		r.logger.Error(
			"activity list command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	// Report table
	tsFromStr, tsToStr := formatTimestamp(tsFrom), formatTimestamp(tsTo)

	htmlBuilder := html.NewBuilder("Таблица ккал активности")
	accordion := html.NewAccordion("accordionActivity")

	// Table
	tbl := html.NewTable([]string{"Дата", "Вес"})

	xlabels := make([]string, 0, len(lst))
	data := make([]float64, 0, len(lst))
	for _, a := range lst {
		tbl.AddRow(
			html.NewTr(nil).
				AddTd(html.NewTd(html.NewS(formatTimestamp(a.Timestamp)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.1f", a.ActiveCal)), nil)),
		)
		xlabels = append(xlabels, formatTimestamp(a.Timestamp))
		data = append(data, a.ActiveCal)
	}

	accordion.AddItem(
		html.HewAccordionItem(
			"tbl",
			fmt.Sprintf("Таблица ккал активности за %s - %s", tsFromStr, tsToStr),
			tbl))

	// Chart
	chart := html.NewCanvas("chart")
	accordion.AddItem(
		html.HewAccordionItem(
			"graph",
			fmt.Sprintf("График ккал активности за %s - %s", tsFromStr, tsToStr),
			chart))

	chartSnip, err := GetChartSnippet(&ChartData{
		ElemID:  "chart",
		XLabels: xlabels,
		Type:    "line",
		Datasets: []ChartDataset{
			{
				Data:  data,
				Label: "ККал",
				Color: ChartColorBlue,
			},
		},
	})
	if err != nil {
		r.logger.Error(
			"activity list command chart error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			accordion,
		),
		html.NewScript(_jsBootstrapURL),
		html.NewScript(_jsChartURL),
		html.NewS(chartSnip),
	)

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("activity_%s_%s.html", tsFromStr, tsToStr),
	})
}

func (r *CmdProcessor) activityDelCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid activity del command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Parse timestamp
	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid activity del command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteActivity(ctx, userID, ts); err != nil {
		r.logger.Error(
			"activity del command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}
