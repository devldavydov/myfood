package cmdproc

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/devldavydov/myfood/internal/common/messages"
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
)

func (r *CmdProcessor) processUserSettings(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid user settings command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "set":
		resp = r.userSettingsSetCommand(cmdParts[1:], userID)
	case "get":
		resp = r.userSettingsGetCommand(userID)
	case "st":
		resp = r.userSettingsSetTemplateCommand(userID)
	default:
		r.logger.Error(
			"invalid user settings command",
			zap.String("reason", "unknown command"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		resp = NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	return resp
}

func (r *CmdProcessor) userSettingsSetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid user settings set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
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
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	defaultActiveCal, err := strconv.ParseFloat(cmdParts[1], 64)
	if err != nil {
		r.logger.Error(
			"invalid user settings set command",
			zap.String("reason", "default active cal format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetUserSettings(
		ctx,
		userID,
		&storage.UserSettings{
			CalLimit:         calLimit,
			DefaultActiveCal: defaultActiveCal,
		}); err != nil {
		if errors.Is(err, storage.ErrUserSettingsInvalid) {
			return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
		}

		r.logger.Error(
			"user settings set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) userSettingsGetCommand(userID int64) []CmdResponse {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	stgs, err := r.stg.GetUserSettings(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserSettingsNotFound) {
			return NewSingleCmdResponse(messages.MsgErrUserSettingsNotFound)
		}

		r.logger.Error(
			"user setting get command DB error",
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("УБМ: %.2f\nАктивные ккал по-умолчанию: %.2f", stgs.CalLimit, stgs.DefaultActiveCal))
}

func (r *CmdProcessor) userSettingsSetTemplateCommand(userID int64) []CmdResponse {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	stgs, err := r.stg.GetUserSettings(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserSettingsNotFound) {
			return NewSingleCmdResponse(messages.MsgErrUserSettingsNotFound)
		}

		r.logger.Error(
			"user setting set template command DB error",
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	usSetTemplate := fmt.Sprintf(
		"us,set,%.2f,%.2f",
		stgs.CalLimit,
		stgs.DefaultActiveCal,
	)

	return NewSingleCmdResponse(usSetTemplate)
}
