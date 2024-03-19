package myfoodbot

import (
	"context"
	"encoding/csv"
	"errors"
	"strconv"
	"strings"

	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
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
	return c.Send("Help! I need somebody Help!")
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
	case "del":
		return r.weightDelCommand(c, cmdParts[1:], userID)
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
	if err := r.stg.CreateWeight(context.Background(), userID, &storage.Weight{Timestamp: ts, Value: val}); err != nil {
		r.logger.Error(
			"weight add command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		if errors.Is(err, storage.ErrWeightAlreadyExists) {
			return c.Send(msgErrWeightAlreadyExists)
		}

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
	if err := r.stg.DeleteWeight(context.Background(), userID, ts); err != nil {
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
