package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) processJournal(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid journal command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	switch cmdParts[0] {
	case "set":
		return r.journalSetCommand(c, cmdParts[1:], userID)
	case "del":
		return r.journalDelCommand(c, cmdParts[1:], userID)
	case "rd":
		return r.journalReportDayCommand(c, cmdParts[1:], userID)
	case "rw":
		return r.journalReportWeek(c, cmdParts[1:], userID)
	}

	r.logger.Error(
		"invalid journal command",
		zap.String("reason", "unknown command"),
		zap.Strings("command", cmdParts),
		zap.Int64("userid", userID),
	)
	return c.Send(msgErrInvalidCommand)
}

func (r *CmdProcessor) journalSetCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 4 {
		r.logger.Error(
			"invalid journal set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	jrnl := &storage.Journal{
		Meal:    storage.NewMealFromString(cmdParts[1]),
		FoodKey: cmdParts[2],
	}

	// Parse timestamp
	ts, err := parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid journal set command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}
	jrnl.Timestamp = ts

	// Parse weight
	weight, err := strconv.ParseFloat(cmdParts[3], 64)
	if err != nil {
		r.logger.Error(
			"invalid weight set command",
			zap.String("reason", "val format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}
	jrnl.FoodWeight = weight

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.SetJournal(ctx, userID, jrnl); err != nil {
		if errors.Is(err, storage.ErrJournalInvalid) {
			return c.Send(msgErrInvalidCommand)
		}

		if errors.Is(err, storage.ErrJournalInvalidFood) {
			return c.Send(msgErrFoodNotFound)
		}

		r.logger.Error(
			"journal set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *CmdProcessor) journalDelCommand(c tele.Context, cmdParts []string, userID int64) error {
	return nil
}

func (r *CmdProcessor) journalReportDayCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid journal rdm command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	tsStr := cmdParts[0]
	ts, err := parseTimestamp(tsStr)
	if err != nil {
		r.logger.Error(
			"invalid journal rd command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// List from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	lst, err := r.stg.GetJournalForPeriod(ctx, userID, ts, ts)
	if err != nil {
		if errors.Is(err, storage.ErrJournalReportEmpty) {
			return c.Send(msgErrEmptyList)
		}

		r.logger.Error(
			"journal rd command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	// Prepare report
	var sb strings.Builder

	sb.WriteString("<html>")
	sb.WriteString("<table border=\"1\" width=\"100%\">")
	sb.WriteString(`<tr>
		<th>Наименование</th>
		<th>Вес</th>
		<th>ККал</th>
		<th>Белки</th>
		<th>Жиры</th>
		<th>Углеводы</th>
	</tr>`)

	var totalCal, totalProt, totalFat, totalCarb float64
	lastMeal := ""
	for _, j := range lst {
		ml := j.Meal.ToString()
		if ml != lastMeal {
			sb.WriteString(fmt.Sprintf(`<tr><td colspan="6" align="center"><b>%s</b><tr>`, ml))
			lastMeal = ml
		}

		foobLbl := j.FoodName
		if j.FoodBrand != "" {
			foobLbl = fmt.Sprintf("%s (%s)", foobLbl, j.FoodBrand)
		}
		sb.WriteString(
			fmt.Sprintf("<tr><td>%s</td><td>%.1f</td><td>%.1f</td><td>%.1f</td><td>%.1f</td><td>%.1f</td></tr>",
				foobLbl,
				j.FoodWeight,
				j.Cal,
				j.Prot,
				j.Fat,
				j.Carb))

		totalCal += j.Cal
		totalProt += j.Prot
		totalFat += j.Fat
		totalCarb += j.Carb
	}

	totalPFC := totalProt + totalFat + totalCarb

	sb.WriteString("</table>")
	sb.WriteString(fmt.Sprintf("<p><b>Всего, ккал: </b>%.1f</p>", totalCal))

	if totalPFC != 0 {
		sb.WriteString(fmt.Sprintf("<p><b>Всего, Б: </b>%.1f (%.1f %%)</p>", totalProt, totalProt/totalPFC*100))
		sb.WriteString(fmt.Sprintf("<p><b>Всего, Ж: </b>%.1f (%.1f %%)</p>", totalFat, totalFat/totalPFC*100))
		sb.WriteString(fmt.Sprintf("<p><b>Всего, У: </b>%.1f (%.1f %%)</p>", totalCarb, totalCarb/totalPFC*100))
	} else {
		sb.WriteString(fmt.Sprintf("<p><b>Всего, Б: </b>%.1f</p>", totalProt))
		sb.WriteString(fmt.Sprintf("<p><b>Всего, Ж: </b>%.1f</p>", totalFat))
		sb.WriteString(fmt.Sprintf("<p><b>Всего, У: </b>%.1f</p>", totalCarb))
	}

	return c.Send(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(sb.String())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("journal_%s.html", tsStr),
	})
}

func (r *CmdProcessor) journalReportWeek(c tele.Context, cmdParts []string, userID int64) error {
	return nil
}
