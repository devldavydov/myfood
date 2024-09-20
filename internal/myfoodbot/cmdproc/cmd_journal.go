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

	"github.com/devldavydov/myfood/internal/common/html"
	"github.com/devldavydov/myfood/internal/common/messages"
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) processJournal(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid journal command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "set":
		resp = r.journalSetCommand(cmdParts[1:], userID)
	case "sb":
		resp = r.journalSetBundleCommand(cmdParts[1:], userID)
	case "del":
		resp = r.journalDelCommand(cmdParts[1:], userID)
	case "dm":
		resp = r.journalDelMealCommand(cmdParts[1:], userID)
	case "cp":
		resp = r.journalCopyCommand(cmdParts[1:], userID)
	case "rm":
		resp = r.journalReportMealCommand(cmdParts[1:], userID)
	case "rd":
		resp = r.journalReportDayCommand(cmdParts[1:], userID)
	case "rw":
		resp = r.journalReportWeekCommand(cmdParts[1:], userID)
	case "rr":
		resp = r.journalReportRangeCommand(cmdParts[1:], userID)
	case "tm":
		resp = r.journalTemplateMealCommand(cmdParts[1:], userID)
	case "fa":
		resp = r.journalFoodAvgWeightCommand(cmdParts[1:], userID)
	default:
		r.logger.Error(
			"invalid journal command",
			zap.String("reason", "unknown command"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		resp = NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	return resp
}

func (r *CmdProcessor) journalSetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 4 {
		r.logger.Error(
			"invalid journal set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	jrnl := &storage.Journal{
		Meal:    storage.NewMealFromString(cmdParts[1]),
		FoodKey: cmdParts[2],
	}

	// Parse timestamp
	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid journal set command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
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
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}
	jrnl.FoodWeight = weight

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetJournal(ctx, userID, jrnl); err != nil {
		if errors.Is(err, storage.ErrJournalInvalid) {
			return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
		}

		if errors.Is(err, storage.ErrJournalInvalidFood) {
			return NewSingleCmdResponse(messages.MsgErrFoodNotFound)
		}

		r.logger.Error(
			"journal set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) journalSetBundleCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 3 {
		r.logger.Error(
			"invalid journal set bundle command",
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
			"invalid journal set bundle command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	meal := storage.NewMealFromString(cmdParts[1])
	bndlKey := cmdParts[2]

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetJournalBundle(ctx, userID, ts, meal, bndlKey); err != nil {
		if errors.Is(err, storage.ErrJournalInvalidFood) {
			return NewSingleCmdResponse(messages.MsgErrFoodNotFound)
		}
		if errors.Is(err, storage.ErrBundleNotFound) {
			return NewSingleCmdResponse(messages.MsgErrBundleNotFound)
		}

		r.logger.Error(
			"journal set bundle command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) journalDelCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 3 {
		r.logger.Error(
			"invalid journal del command",
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
			"invalid journal del command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteJournal(ctx, userID, ts, storage.NewMealFromString(cmdParts[1]), cmdParts[2]); err != nil {
		r.logger.Error(
			"journal del command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) journalDelMealCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid journal dm command",
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
			"invalid journal dm command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteJournalMeal(ctx, userID, ts, storage.NewMealFromString(cmdParts[1])); err != nil {
		r.logger.Error(
			"journal dm command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(messages.MsgOK)
}

func (r *CmdProcessor) journalCopyCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 4 {
		r.logger.Error(
			"invalid journal copy command",
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
			"invalid journal copy command",
			zap.String("reason", "ts from format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	tsTo, err := r.parseTimestamp(cmdParts[2])
	if err != nil {
		r.logger.Error(
			"invalid journal copy command",
			zap.String("reason", "ts to format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	cnt, err := r.stg.CopyJournal(ctx,
		userID,
		tsFrom,
		storage.NewMealFromString(cmdParts[1]),
		tsTo,
		storage.NewMealFromString(cmdParts[3]))

	if err != nil {
		if errors.Is(err, storage.ErrCopyToNotEmpty) {
			return NewSingleCmdResponse(messages.MsgErrJournalCopy)
		}

		r.logger.Error(
			"journal copy command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf(messages.MsgJournalCopied, cnt))
}

func (r *CmdProcessor) journalReportMealCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid journal rm command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid journal rm command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Get list from DB and user settings
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout*2)
	defer cancel()

	rep, err := r.stg.GetJournalMealReport(ctx, userID, ts, storage.NewMealFromString(cmdParts[1]))
	if err != nil {
		if errors.Is(err, storage.ErrJournalMealReportEmpty) {
			return NewSingleCmdResponse(messages.MsgErrEmptyList)
		}

		r.logger.Error(
			"journal rm command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	var us *storage.UserSettings
	us, err = r.stg.GetUserSettings(ctx, userID)
	if err != nil {
		if !errors.Is(err, storage.ErrUserSettingsNotFound) {
			r.logger.Error(
				"journal rm command DB error for user settings",
				zap.Strings("command", cmdParts),
				zap.Int64("userid", userID),
				zap.Error(err),
			)

			return NewSingleCmdResponse(messages.MsgErrInternal)
		}
	}

	var sb strings.Builder
	for _, item := range rep.Items {
		foodLbl := item.FoodName
		if item.FoodBrand != "" {
			foodLbl += " - " + item.FoodBrand
		}
		sb.WriteString(fmt.Sprintf("<b>%s [%s]</b>:\n", foodLbl, item.FoodKey))
		sb.WriteString(fmt.Sprintf("%.1f г., %.2f ккал\n", item.FoodWeight, item.Cal))
	}
	sb.WriteString(fmt.Sprintf("\n<b>Всего (прием пищи), ккал:</b> %.2f", rep.ConsumedMealCal))
	if us != nil {
		sb.WriteString(fmt.Sprintf("\n<b>Всего (день), ккал:</b> %.2f (%+.2f)", rep.ConsumedDayCal, us.CalLimit-rep.ConsumedDayCal))
	} else {
		sb.WriteString(fmt.Sprintf("\n<b>Всего (день), ккал:</b> %.2f", rep.ConsumedDayCal))
	}

	return NewSingleCmdResponse(sb.String(), optsHTML)
}

func (r *CmdProcessor) journalReportDayCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid journal rd command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid journal rd command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}
	tsStr := formatTimestamp(ts)

	// Get list from DB, user settings and activity
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout*2)
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

			return NewSingleCmdResponse(messages.MsgErrInternal)
		}
	}

	var ua *storage.Activity
	ua, err = r.stg.GetActivity(ctx, userID, ts)
	if err != nil {
		if !errors.Is(err, storage.ErrActivityNotFound) {
			r.logger.Error(
				"journal rd command DB error for activity",
				zap.Strings("command", cmdParts),
				zap.Int64("userid", userID),
				zap.Error(err),
			)

			return NewSingleCmdResponse(messages.MsgErrInternal)
		}
	}

	lst, err := r.stg.GetJournalReport(ctx, userID, ts, ts)
	if err != nil {
		if errors.Is(err, storage.ErrJournalReportEmpty) {
			return NewSingleCmdResponse(messages.MsgErrEmptyList)
		}

		r.logger.Error(
			"journal rd command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	// Report table
	htmlBuilder := html.NewBuilder("Журнал приема пищи")
	tbl := html.NewTable([]string{
		"Наименование", "Вес", "ККал", "Белки", "Жиры", "Углеводы",
	})

	var totalCal, totalProt, totalFat, totalCarb float64
	var subTotalCal, subTotalProt, subTotalFat, subTotalCarb float64
	lastMeal := storage.Meal(-1)
	for i := 0; i < len(lst); i++ {
		j := lst[i]

		// Add meal divider
		if j.Meal != lastMeal {
			tbl.AddRow(
				html.NewTr(html.Attrs{"class": "table-active"}).
					AddTd(html.NewTd(
						html.NewB(j.Meal.ToString(), nil),
						html.Attrs{"colspan": "6", "align": "center"},
					)),
			)
			lastMeal = j.Meal
		}

		// Add meal rows
		foodLbl := j.FoodName
		if j.FoodBrand != "" {
			foodLbl = fmt.Sprintf("%s - %s", foodLbl, j.FoodBrand)
		}
		foodLbl = fmt.Sprintf("%s [%s]", foodLbl, j.FoodKey)

		tbl.AddRow(
			html.NewTr(nil).
				AddTd(html.NewTd(html.NewS(foodLbl), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.1f", j.FoodWeight)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.Cal)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.Prot)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.Fat)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.Carb)), nil)))

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
			tbl.AddRow(
				html.NewTr(nil).
					AddTd(html.NewTd(html.NewB("Всего", nil), html.Attrs{"align": "right", "colspan": "2"})).
					AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", subTotalCal)), nil)).
					AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", subTotalProt)), nil)).
					AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", subTotalFat)), nil)).
					AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", subTotalCarb)), nil)))

			subTotalCal, subTotalProt, subTotalFat, subTotalCarb = 0, 0, 0, 0
		}
	}

	// Footer
	totalPFC := totalProt + totalFat + totalCarb

	tbl.
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Всего, ккал: ", nil),
						html.NewS(fmt.Sprintf("%.2f", totalCal)),
					),
					html.Attrs{"colspan": "6"})))

	if us != nil {
		activeCal := us.DefaultActiveCal
		if ua != nil {
			activeCal = ua.ActiveCal
		}

		tbl.AddFooterElement(html.NewTr(nil).
			AddTd(html.NewTd(
				html.NewSpan(
					html.NewB("УБМ, ккал: ", nil),
					html.NewS(fmt.Sprintf("%.2f", us.CalLimit)),
					html.NewNbsp(),
					html.NewB("Активность, ккал: ", nil),
					html.NewS(fmt.Sprintf("%.2f", activeCal)),
					html.NewNbsp(),
					html.NewB("Разница, ккал: ", nil),
					calDiffSnippet2(us.CalLimit+activeCal-totalCal),
				),
				html.Attrs{"colspan": "6"})))
	}

	tbl.
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Всего, Б: ", nil),
						pfcSnippet(totalProt, totalPFC),
					),
					html.Attrs{"colspan": "6"}))).
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Всего, Ж: ", nil),
						pfcSnippet(totalFat, totalPFC),
					),
					html.Attrs{"colspan": "6"}))).
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Всего, У: ", nil),
						pfcSnippet(totalCarb, totalPFC),
					),
					html.Attrs{"colspan": "6"})))

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			html.NewH(
				fmt.Sprintf("Журнал приема пищи за %s", tsStr),
				5,
				html.Attrs{"align": "center"},
			),
			tbl,
		),
	)

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("report_%s.html", tsStr),
	})
}

func (r *CmdProcessor) journalReportWeekCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid journal rw command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	tsStart, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid journal rw command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	tsStart = getStartOfWeek(tsStart)
	tsStartUnix := tsStart.Unix()
	tsStartStr := formatTimestamp(tsStart)

	tsRangeStr := make([]string, 7)
	for i := 0; i < 7; i++ {
		tsRangeStr[i] = formatTimestamp(tsStart.Add(time.Duration(i) * 24 * time.Hour))
	}
	tsEnd := tsStart.Add(6 * 24 * time.Hour)
	tsEndStr := formatTimestamp(tsEnd)

	// Get list from DB and user settings
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout*2)
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

			return NewSingleCmdResponse(messages.MsgErrInternal)
		}
	}

	lst, err := r.stg.GetJournalStats(ctx, userID, tsStart, tsEnd)
	if err != nil {
		if errors.Is(err, storage.ErrJournalStatsEmpty) {
			return NewSingleCmdResponse(messages.MsgErrEmptyList)
		}

		r.logger.Error(
			"journal rw command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	// Stat table
	htmlBuilder := html.NewBuilder("Статистика приема пищи")
	accordion := html.NewAccordion("accordionStats")

	// Table
	tbl := html.NewTable([]string{
		"Дата", "Итого, ккал", "Итого, белки", "Итого, жиры", "Итого, углеводы",
	})

	var totalCal, totalProt, totalFat, totalCarb float64
	dataRange := make([]float64, 7)
	for _, j := range lst {
		tbl.AddRow(
			html.NewTr(nil).
				AddTd(html.NewTd(html.NewS(formatTimestamp(j.Timestamp)), nil)).
				AddTd(html.NewTd(calDiffSnippet(us, j.TotalCal), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.TotalProt)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.TotalFat)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.TotalCarb)), nil)))

		totalCal += j.TotalCal
		totalProt += j.TotalProt
		totalFat += j.TotalFat
		totalCarb += j.TotalCarb

		dataRange[(j.Timestamp.Unix()-tsStartUnix)/24/3600%7] = j.TotalCal
	}

	// Footer
	lLst := float64(len(lst))
	avgCal, avgProt, avgFat, avgCarb := totalCal/lLst, totalProt/lLst, totalFat/lLst, totalCarb/lLst
	totalAvgPFC := avgProt + avgFat + avgCarb

	tbl.
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Среднее, ккал: ", nil),
						html.NewS(fmt.Sprintf("%.2f", avgCal)),
					),
					html.Attrs{"colspan": "5"}))).
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Среднее, Б: ", nil),
						pfcSnippet(avgProt, totalAvgPFC),
					),
					html.Attrs{"colspan": "5"}))).
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Среднее, Ж: ", nil),
						pfcSnippet(avgFat, totalAvgPFC),
					),
					html.Attrs{"colspan": "5"}))).
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Среднее, У: ", nil),
						pfcSnippet(avgCarb, totalAvgPFC),
					),
					html.Attrs{"colspan": "5"})))

	accordion.AddItem(
		html.HewAccordionItem(
			"tbl",
			fmt.Sprintf("Статистика приема пищи за %s - %s", tsStartStr, tsEndStr),
			tbl))

	// Chart
	chart := html.NewCanvas("chart")
	accordion.AddItem(
		html.HewAccordionItem(
			"graph",
			fmt.Sprintf("График приема пищи за %s - %s", tsStartStr, tsEndStr),
			chart))

	data := &ChardData{
		ElemID:  "chart",
		XLabels: tsRangeStr,
		Data:    dataRange,
		Label:   "ККал",
		Type:    "bar",
	}
	chartSnip, err := GetChartSnippet(data)
	if err != nil {
		r.logger.Error(
			"journal rw command chart error",
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

	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("stats_%s_%s.html", tsStartStr, tsEndStr),
	})
}

func (r *CmdProcessor) journalReportRangeCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid journal rr command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	tsStart, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid journal rr command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	tsEnd, err := r.parseTimestamp(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid journal rr command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Get list from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout*2)
	defer cancel()

	lst, err := r.stg.GetJournalStats(ctx, userID, tsStart, tsEnd)
	if err != nil {
		if errors.Is(err, storage.ErrJournalStatsEmpty) {
			return NewSingleCmdResponse(messages.MsgErrEmptyList)
		}

		r.logger.Error(
			"journal rr command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	// Get chart data
	xlabels := make([]string, 0, len(lst))
	data := make([]float64, 0, len(lst))
	for _, w := range lst {
		xlabels = append(xlabels, formatTimestamp(w.Timestamp))
		data = append(data, w.TotalCal)
	}

	// HTML report
	tsStartStr := formatTimestamp(tsStart)
	tsEndStr := formatTimestamp(tsEnd)

	// Chart
	chartSnip, err := GetChartSnippet(&ChardData{
		ElemID:  "chart",
		XLabels: xlabels,
		Data:    data,
		Label:   "ККал",
		Type:    "line",
	})
	if err != nil {
		r.logger.Error(
			"weight list command chart error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	// Doc
	htmlBuilder := html.
		NewBuilder("График ККал за период").
		Add(
			html.NewContainer().Add(
				html.NewH(
					fmt.Sprintf(fmt.Sprintf("График ККал за период %s - %s", tsStartStr, tsEndStr)),
					5,
					html.Attrs{"align": "center"},
				),
				html.NewCanvas("chart"),
				html.NewScript(_jsBootstrapURL),
				html.NewScript(_jsChartURL),
				html.NewS(chartSnip),
			),
		)

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("report_range_%s_%s.html", tsStartStr, tsEndStr),
	})
}

func (r *CmdProcessor) journalTemplateMealCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid journal tm command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid journal tm command",
			zap.String("reason", "ts format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	// Get list from DB and user settings
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout*2)
	defer cancel()

	mealStr := cmdParts[1]
	rep, err := r.stg.GetJournalMealReport(ctx, userID, ts, storage.NewMealFromString(mealStr))
	if err != nil {
		if errors.Is(err, storage.ErrJournalMealReportEmpty) {
			return NewSingleCmdResponse(messages.MsgErrEmptyList)
		}

		r.logger.Error(
			"journal tm command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	tsStr := formatTimestamp(ts)
	resp := make([]CmdResponse, 0)

	resp = append(resp, NewCmdResponse("<b>Изменение еды</b>", optsHTML))
	for _, item := range rep.Items {
		resp = append(resp, NewCmdResponse(
			fmt.Sprintf("j,set,%s,%s,%s,%.1f", tsStr, mealStr, item.FoodKey, item.FoodWeight),
		))
	}
	resp = append(resp, NewCmdResponse("<b>Удаление еды</b>", optsHTML))
	for _, item := range rep.Items {
		resp = append(resp, NewCmdResponse(
			fmt.Sprintf("j,del,%s,%s,%s", tsStr, mealStr, item.FoodKey),
		))
	}

	return resp
}

func (r *CmdProcessor) journalFoodAvgWeightCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid journal fa command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return NewSingleCmdResponse(messages.MsgErrInvalidCommand)
	}

	tsTo := time.Now()
	tsFrom := tsTo.AddDate(-1, 0, 0)

	// Get data from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout*2)
	defer cancel()

	avgW, err := r.stg.GetJournalFoodAvgWeight(ctx, userID, tsFrom, tsTo, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrJournalInvalidFood) {
			return NewSingleCmdResponse(messages.MsgErrFoodNotFound)
		}

		r.logger.Error(
			"journal fa command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(messages.MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("Средний вес прима пищи за год: %.1fг.", avgW))
}

func calDiffSnippet(us *storage.UserSettings, cal float64) html.IELement {
	if us == nil {
		return html.NewS(fmt.Sprintf("%.2f", cal))
	} else {
		diff := us.CalLimit - cal
		switch {
		case diff < 0 && math.Abs(diff) > 0.01:
			return html.NewSpan(
				html.NewS(fmt.Sprintf("%.2f (", cal)),
				html.NewB(fmt.Sprintf("%+.2f", diff), html.Attrs{"class": "text-danger"}),
				html.NewS(")"),
			)
		case diff > 0 && math.Abs(diff) > 0.01:
			return html.NewSpan(
				html.NewS(fmt.Sprintf("%.2f (", cal)),
				html.NewB(fmt.Sprintf("%+.2f", diff), html.Attrs{"class": "text-success"}),
				html.NewS(")"),
			)
		default:
			return html.NewS(fmt.Sprintf("%.2f", cal))
		}
	}
}

func calDiffSnippet2(diff float64) html.IELement {
	switch {
	case diff < 0 && math.Abs(diff) > 0.01:
		return html.NewSpan(
			html.NewB(fmt.Sprintf("%+.2f", diff), html.Attrs{"class": "text-danger"}),
		)
	case diff >= 0 && math.Abs(diff) > 0.01:
		return html.NewSpan(
			html.NewB(fmt.Sprintf("%+.2f", diff), html.Attrs{"class": "text-success"}),
		)
	default:
		return html.NewS(fmt.Sprintf("%.2f", diff))
	}
}

func pfcSnippet(val, totalVal float64) html.IELement {
	var s string

	if totalVal == 0 {
		s = fmt.Sprintf("%.2f", val)
	} else {
		s = fmt.Sprintf("%.2f (%.2f%%)", val, val/totalVal*100)
	}

	return html.NewS(s)
}
