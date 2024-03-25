package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

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
	if len(cmdParts) != 3 {
		r.logger.Error(
			"invalid journal del command",
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
			"invalid journal del command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteJournal(ctx, userID, ts, storage.NewMealFromString(cmdParts[1]), cmdParts[2]); err != nil {
		r.logger.Error(
			"weight journal command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(msgOK)
}

func (r *CmdProcessor) journalReportDayCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid journal rd command",
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

	// Get list from DB and user settings
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout*2)
	defer cancel()

	var us *storage.UserSettings
	us, err = r.stg.GetUserSettings(ctx, userID)
	if err != nil {
		if !errors.Is(err, storage.ErrUserSettingsNotFound) {
			r.logger.Error(
				"journal rd command DB error for user settings",
				zap.Strings("command", cmdParts),
				zap.Int64("userid", userID),
				zap.Error(err),
			)

			return c.Send(msgErrInternal)
		}
	}

	lst, err := r.stg.GetJournalReport(ctx, userID, ts, ts)
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

	// Report table
	var sb strings.Builder

	sb.WriteString("<html>")
	sb.WriteString(fmt.Sprintf(`<link href="%s" rel="stylesheet">`, _styleURL))
	sb.WriteString(`<div class="container">`)
	sb.WriteString(fmt.Sprintf(`<h5 align="center">Журнал приема пищи за %s</h5>`, tsStr))
	sb.WriteString(`<table class="table table-bordered table-hover">`)
	sb.WriteString(`<thead class="table-light">
		<tr>
			<th>Наименование</th>
			<th>Вес</th>
			<th>ККал</th>
			<th>Белки</th>
			<th>Жиры</th>
			<th>Углеводы</th>
		</tr>
	</thead>`)

	// Body
	sb.WriteString("<tbody>")
	var totalCal, totalProt, totalFat, totalCarb float64
	var subTotalCal, subTotalProt, subTotalFat, subTotalCarb float64
	lastMeal := storage.Meal(-1)
	for i := 0; i < len(lst); i++ {
		j := lst[i]

		// Add meal divider
		if j.Meal != lastMeal {
			sb.WriteString(fmt.Sprintf(`<tr><td colspan="6" align="center"><b>%s</b><tr>`, j.Meal.ToString()))
			lastMeal = j.Meal
		}

		// Add meal rows
		foodLbl := j.FoodName
		if j.FoodBrand != "" {
			foodLbl = fmt.Sprintf("%s - %s", foodLbl, j.FoodBrand)
		}
		foodLbl = fmt.Sprintf("%s [%s]", foodLbl, j.FoodKey)

		sb.WriteString(
			fmt.Sprintf("<tr><td>%s</td><td>%.1f</td><td>%.2f</td><td>%.2f</td><td>%.2f</td><td>%.2f</td></tr>",
				foodLbl, j.FoodWeight, j.Cal, j.Prot, j.Fat, j.Carb))

		totalCal += j.Cal
		totalProt += j.Prot
		totalFat += j.Fat
		totalCarb += j.Carb

		subTotalCal += j.Cal
		subTotalProt += j.Prot
		subTotalFat += j.Fat
		subTotalCarb += j.Carb

		// Add subtotal row
		if i == len(lst)-1 || lst[i+1].Meal != j.Meal {
			sb.WriteString(fmt.Sprintf(`<tr>
			<td align="right" colspan="2"><b>Всего</b></td>
			<td>%.2f</td>
			<td>%.2f</td>
			<td>%.2f</td>
			<td>%.2f</td>
			</tr>`, subTotalCal, subTotalProt, subTotalFat, subTotalCarb))
			subTotalCal, subTotalProt, subTotalFat, subTotalCarb = 0, 0, 0, 0
		}
	}
	sb.WriteString("</tbody>")

	// Footer
	sb.WriteString("<tfoot>")
	sb.WriteString(fmt.Sprintf(`<tr><td colspan="6"><b>Всего, ккал: </b>%s</td></tr>`, calDiffSnippet(us, totalCal)))

	totalPFC := totalProt + totalFat + totalCarb
	sb.WriteString(fmt.Sprintf(`<tr><td colspan="6"><b>Всего, Б: </b>%s</td></tr>`, pfcSnippet(totalProt, totalPFC)))
	sb.WriteString(fmt.Sprintf(`<tr><td colspan="6"><b>Всего, Ж: </b>%s</td></tr>`, pfcSnippet(totalFat, totalPFC)))
	sb.WriteString(fmt.Sprintf(`<tr><td colspan="6"><b>Всего, У: </b>%s</td></tr>`, pfcSnippet(totalCarb, totalPFC)))
	sb.WriteString("</tfoot>")

	// End
	sb.WriteString("</table></div></html>")

	return c.Send(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(sb.String())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("report_%s.html", tsStr),
	})
}

func (r *CmdProcessor) journalReportWeek(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid journal rw command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return c.Send(msgErrInvalidCommand)
	}

	tsStartStr := cmdParts[0]
	tsStart, err := parseTimestamp(tsStartStr)
	if err != nil {
		r.logger.Error(
			"invalid journal rw command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	if !isStartOfWeek(tsStart) {
		return c.Send(msgErrJournalNotStartOfWeek)
	}

	tsEnd := time.Unix(tsStart, 0).Add(6 * 24 * time.Hour).Unix()
	tsEndStr := formatTimestamp(tsEnd)

	// Get list from DB and user settings
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout*2)
	defer cancel()

	var us *storage.UserSettings
	us, err = r.stg.GetUserSettings(ctx, userID)
	if err != nil {
		if !errors.Is(err, storage.ErrUserSettingsNotFound) {
			r.logger.Error(
				"journal rw command DB error for user settings",
				zap.Strings("command", cmdParts),
				zap.Int64("userid", userID),
				zap.Error(err),
			)

			return c.Send(msgErrInternal)
		}
	}

	lst, err := r.stg.GetJournalStats(ctx, userID, tsStart, tsEnd)
	if err != nil {
		if errors.Is(err, storage.ErrJournalStatsEmpty) {
			return c.Send(msgErrEmptyList)
		}

		r.logger.Error(
			"journal rw command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	// Stat table
	var sb strings.Builder

	sb.WriteString("<html>")
	sb.WriteString(fmt.Sprintf(`<link href="%s" rel="stylesheet">`, _styleURL))
	sb.WriteString(`<div class="container">`)
	sb.WriteString(fmt.Sprintf(`<h5 align="center">Статистика приема пищи за %s - %s</h5>`, tsStartStr, tsEndStr))
	sb.WriteString(`<table class="table table-bordered table-hover">`)
	sb.WriteString(`<thead class="table-light">
		<tr>
			<th>Дата</th>
			<th>Итого, ккал</th>
			<th>Итого, белки</th>
			<th>Итого, жиры</th>
			<th>Итого углеводы</th>
		</tr>
	</thead>`)

	// Body
	sb.WriteString("<tbody>")
	var totalCal, totalProt, totalFat, totalCarb float64
	for _, j := range lst {
		sb.WriteString(
			fmt.Sprintf(
				"<tr><td>%s</td><td>%s</td><td>%.2f</td><td>%.2f</td><td>%.2f</td></tr>",
				formatTimestamp(j.Timestamp),
				calDiffSnippet(us, j.TotalCal),
				j.TotalProt,
				j.TotalFat,
				j.TotalCarb))

		totalCal += j.TotalCal
		totalProt += j.TotalProt
		totalFat += j.TotalFat
		totalCarb += j.TotalCarb
	}
	sb.WriteString("</tbody>")

	lLst := float64(len(lst))
	avgCal, avgProt, avgFat, avgCarb := totalCal/lLst, totalProt/lLst, totalFat/lLst, totalCarb/lLst

	// Footer
	sb.WriteString("<tfoot>")
	sb.WriteString(fmt.Sprintf(`<tr><td colspan="5"><b>Среднее, ккал: </b>%.2f</td></tr>`, avgCal))

	totalAvgPFC := avgProt + avgFat + avgCarb
	sb.WriteString(fmt.Sprintf(`<tr><td colspan="5"><b>Среднее, Б: </b>%s</td></tr>`, pfcSnippet(avgProt, totalAvgPFC)))
	sb.WriteString(fmt.Sprintf(`<tr><td colspan="5"><b>Среднее, Ж: </b>%s</td></tr>`, pfcSnippet(avgFat, totalAvgPFC)))
	sb.WriteString(fmt.Sprintf(`<tr><td colspan="5"><b>Среднее, У: </b>%s</td></tr>`, pfcSnippet(avgCarb, totalAvgPFC)))
	sb.WriteString("</tfoot>")

	// End
	sb.WriteString("</table></div></html>")

	return c.Send(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(sb.String())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("stats_%s_%s.html", tsStartStr, tsEndStr),
	})
}

func calDiffSnippet(us *storage.UserSettings, cal float64) string {
	if us == nil {
		return fmt.Sprintf("%.2f", cal)
	} else {
		diff := us.CalLimit - cal
		switch {
		case diff < 0 && math.Abs(diff) > 0.01:
			return fmt.Sprintf(`%.2f (<b class="text-danger">%+.2f</b>)`, cal, diff)
		case diff > 0 && math.Abs(diff) > 0.01:
			return fmt.Sprintf(`%.2f (<b class="text-success">%+.2f</b>)`, cal, diff)
		default:
			return fmt.Sprintf("%.2f", cal)
		}
	}
}

func pfcSnippet(val, totalVal float64) string {
	if totalVal == 0 {
		return fmt.Sprintf("%.2f", val)
	}

	return fmt.Sprintf("%.2f (%.2f%%)", val, val/totalVal*100)
}