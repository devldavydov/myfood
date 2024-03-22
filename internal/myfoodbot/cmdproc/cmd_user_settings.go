package cmdproc

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) processUserSettings(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid user settings command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	switch cmdParts[0] {
	case "set":
		return r.userSettingsSetCommand(c, cmdParts[1:], userID)
	case "get":
		return r.userSettingsGetCommand(c, userID)
	}

	r.logger.Error(
		"invalid user settings command",
		zap.String("reason", "unknown command"),
		zap.Strings("command", cmdParts),
		zap.Int64("userid", userID),
	)
	return c.Send(msgErrInvalidCommand)
}

func (r *CmdProcessor) userSettingsSetCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid user settings set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// parse
	calLimit, err := strconv.ParseFloat(cmdParts[0], 64)
	if err != nil {
		r.logger.Error(
			"invalid user settings set command",
			zap.String("reason", "cal limit format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.SetUserSettings(ctx, userID, &storage.UserSettings{CalLimit: calLimit}); err != nil {
		if errors.Is(err, storage.ErrUserSettingsInvalid) {
			return c.Send(msgErrInvalidCommand)
		}

		r.logger.Error(
			"user settings set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *CmdProcessor) userSettingsGetCommand(c tele.Context, userID int64) error {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	stgs, err := r.stg.GetUserSettings(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserSettingsNotFound) {
			return c.Send(msgErrUserSettingsNotFound)
		}

		r.logger.Error(
			"user setting get command DB error",
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(fmt.Sprintf("Лимит калорий: %.2f", stgs.CalLimit))
}
