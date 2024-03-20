package myfoodbot

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/devldavydov/myfood/internal/myfoodbot/help"
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	tele "gopkg.in/telebot.v3"
)

var (
	_stgOperationTimeout = 10 * time.Second
)

type cmdProcessor struct {
	stg    storage.Storage
	logger *zap.Logger
}

func newCmdProcessor(stg storage.Storage, logger *zap.Logger) *cmdProcessor {
	return &cmdProcessor{stg: stg, logger: logger}
}

func (r *cmdProcessor) process(c tele.Context, cmd string, userID int64) error {
	rd := csv.NewReader(strings.NewReader(cmd))
	cmdParts, err := rd.Read()
	if err != nil {
		r.logger.Error("command parse error", zap.String("command", cmd), zap.Int64("userid", userID), zap.Error(err))
		return c.Send(msgErrInternal)
	}

	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid command",
			zap.String("reason", "empty command"),
			zap.String("command", cmd),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	switch cmdParts[0] {
	case "h":
		return r.helpCommand(c)
	case "w":
		return r.processWeight(c, cmdParts[1:], userID)
	}

	r.logger.Error(
		"invalid command",
		zap.String("reason", "unknown command"),
		zap.String("command", cmd),
		zap.Int64("userid", userID),
	)
	return c.Send(msgErrInvalidCommand)
}

//
// Help command.
//

func (r *cmdProcessor) helpCommand(c tele.Context) error {
	docRd, err := help.GetHelpDocument("help_ru")
	if err != nil {
		r.logger.Error(
			"help command error",
			zap.Int64("userid", c.Sender().ID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}
	return c.Send(&tele.Document{
		File:     tele.FromReader(docRd),
		MIME:     "text/html",
		FileName: "help_ru.html",
	})
}

//
// Weight commands.
//

func (r *cmdProcessor) processWeight(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid weight command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	switch cmdParts[0] {
	case "add":
		return r.weightAddCommand(c, cmdParts[1:], userID)
	case "upd":
		return r.weightUpdCommand(c, cmdParts[1:], userID)
	case "del":
		return r.weightDelCommand(c, cmdParts[1:], userID)
	case "list":
		return r.weightListCommand(c, cmdParts[1:], userID)
	case "graph":
		return r.weightGraphCommand(c, cmdParts[1:], userID)
	}

	r.logger.Error(
		"invalid weight command",
		zap.String("reason", "unknown command"),
		zap.Strings("command", cmdParts),
		zap.Int64("userid", userID),
	)
	return c.Send(msgErrInvalidCommand)
}

func (r *cmdProcessor) weightAddCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid weight add command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse timestamp
	ts, err := parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid weight add command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse val
	val, err := strconv.ParseFloat(cmdParts[1], 64)
	if err != nil {
		r.logger.Error(
			"invalid weight add command",
			zap.String("reason", "val format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.CreateWeight(ctx, userID, &storage.Weight{Timestamp: ts, Value: val}); err != nil {
		if errors.Is(err, storage.ErrWeightAlreadyExists) {
			return c.Send(msgErrWeightAlreadyExists)
		}

		r.logger.Error(
			"weight add command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *cmdProcessor) weightUpdCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid weight upd command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse timestamp
	ts, err := parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid weight upd command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse val
	val, err := strconv.ParseFloat(cmdParts[1], 64)
	if err != nil {
		r.logger.Error(
			"invalid weight upd command",
			zap.String("reason", "val format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.UpdateWeight(ctx, userID, &storage.Weight{Timestamp: ts, Value: val}); err != nil {
		if errors.Is(err, storage.ErrWeightNotFound) {
			return c.Send(msgErrWeightNotFound)
		}

		r.logger.Error(
			"weight upd command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *cmdProcessor) weightDelCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid weight del command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse timestamp
	ts, err := parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid weight del command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteWeight(ctx, userID, ts); err != nil {
		r.logger.Error(
			"weight del command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *cmdProcessor) weightListCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid weight list command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse timestamp
	tsFrom, err := parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid weight list command",
			zap.String("reason", "ts from format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	tsTo, err := parseTimestamp(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid weight list command",
			zap.String("reason", "ts to format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// List from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	lst, err := r.stg.GetWeightList(ctx, userID, tsFrom, tsTo)
	if err != nil {
		if errors.Is(err, storage.ErrWeightEmptyList) {
			return c.Send(msgErrEmptyList)
		}

		r.logger.Error(
			"weight list command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	var sb strings.Builder
	for _, w := range lst {
		sb.WriteString(fmt.Sprintf("%s: %4.1f\n", formatTimestamp(w.Timestamp), w.Value))
	}

	return c.Send(sb.String())
}

func (r *cmdProcessor) weightGraphCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid weight graph command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Parse timestamp
	tsFrom, err := parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid weight graph command",
			zap.String("reason", "ts from format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	tsTo, err := parseTimestamp(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid weight graph command",
			zap.String("reason", "ts to format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// List from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	lst, err := r.stg.GetWeightList(ctx, userID, tsFrom, tsTo)
	if err != nil {
		if errors.Is(err, storage.ErrWeightEmptyList) {
			return c.Send(msgErrEmptyList)
		}

		r.logger.Error(
			"weight graph command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	// Plot graph
	p := plot.New()
	p.Title.Text = "График Веса"
	p.Y.Label.Text = "Вес"

	vals := make(plotter.Values, len(lst))
	nominals := make([]string, len(lst))
	for i := range lst {
		vals[i] = lst[i].Value
		if i == 0 || i == len(lst)-1 {
			nominals[i] = formatTimestamp(lst[i].Timestamp)
			continue
		}
		nominals[i] = ""
	}
	p.NominalX(nominals...)

	bars, err := plotter.NewBarChart(vals, vg.Points(20))
	if err != nil {
		r.logger.Error(
			"weight graph command plot bars",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = plotutil.Color(0)
	p.Add(bars)

	buf := bytes.NewBuffer([]byte{})
	wr, err := p.WriterTo(vg.Points(640), vg.Points(480), "png")
	if err != nil {
		r.logger.Error(
			"weight graph command plot create writer error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	_, err = wr.WriteTo(buf)
	if err != nil {
		r.logger.Error(
			"weight graph command plot writer error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(&tele.Photo{File: tele.FromReader(buf)})
}
