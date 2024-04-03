package cmdproc

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
)

func (r *CmdProcessor) processUserSettings(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid user settings command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return msgErrInvalidCommand, nil
	}

	var what any
	var opts []any

	switch cmdParts[0] {
	case "set":
		what, opts = r.userSettingsSetCommand(cmdParts[1:], userID)
	case "get":
		what, opts = r.userSettingsGetCommand(userID)
	default:
		r.logger.Error(
			"invalid user settings command",
			zap.String("reason", "unknown command"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		what = msgErrInvalidCommand
	}

	return what, opts
}

func (r *CmdProcessor) userSettingsSetCommand(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid user settings set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.SetUserSettings(ctx, userID, &storage.UserSettings{CalLimit: calLimit}); err != nil {
		if errors.Is(err, storage.ErrUserSettingsInvalid) {
			return msgErrInvalidCommand, nil
		}

		r.logger.Error(
			"user settings set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return msgErrInternal, nil
	}

	return msgOK, nil
}

func (r *CmdProcessor) userSettingsGetCommand(userID int64) (any, []any) {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	stgs, err := r.stg.GetUserSettings(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserSettingsNotFound) {
			return msgErrUserSettingsNotFound, nil
		}

		r.logger.Error(
			"user setting get command DB error",
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return msgErrInternal, nil
	}

	return fmt.Sprintf("Лимит калорий: %.2f", stgs.CalLimit), nil
}
