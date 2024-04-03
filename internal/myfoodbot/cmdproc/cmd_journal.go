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
	"github.com/devldavydov/myfood/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (r *CmdProcessor) processJournal(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid journal command",
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
		what, opts = r.journalSetCommand(cmdParts[1:], userID)
	case "del":
		what, opts = r.journalDelCommand(cmdParts[1:], userID)
	case "dm":
		what, opts = r.journalDelMealCommand(cmdParts[1:], userID)
	case "cp":
		what, opts = r.journalCopyCommand(cmdParts[1:], userID)
	case "rm":
		what, opts = r.journalReportMealCommand(cmdParts[1:], userID)
	case "rd":
		what, opts = r.journalReportDayCommand(cmdParts[1:], userID)
	case "rw":
		what, opts = r.journalReportWeek(cmdParts[1:], userID)
	default:
		r.logger.Error(
			"invalid journal command",
			zap.String("reason", "unknown command"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		what = msgErrInvalidCommand
	}

	return what, opts
}

func (r *CmdProcessor) journalSetCommand(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) != 4 {
		r.logger.Error(
			"invalid journal set command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
	}
	jrnl.FoodWeight = weight

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.SetJournal(ctx, userID, jrnl); err != nil {
		if errors.Is(err, storage.ErrJournalInvalid) {
			return msgErrInvalidCommand, nil
		}

		if errors.Is(err, storage.ErrJournalInvalidFood) {
			return msgErrFoodNotFound, nil
		}

		r.logger.Error(
			"journal set command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return msgErrInternal, nil
	}

	return msgOK, nil
}

func (r *CmdProcessor) journalDelCommand(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) != 3 {
		r.logger.Error(
			"invalid journal del command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
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

		return msgErrInternal, nil
	}

	return msgOK, nil
}

func (r *CmdProcessor) journalDelMealCommand(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid journal dm command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteJournalMeal(ctx, userID, ts, storage.NewMealFromString(cmdParts[1])); err != nil {
		r.logger.Error(
			"journal dm command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return msgErrInternal, nil
	}

	return msgOK, nil
}

func (r *CmdProcessor) journalCopyCommand(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) != 4 {
		r.logger.Error(
			"invalid journal copy command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout)
	defer cancel()

	cnt, err := r.stg.CopyJournal(ctx,
		userID,
		tsFrom,
		storage.NewMealFromString(cmdParts[1]),
		tsTo,
		storage.NewMealFromString(cmdParts[3]))

	if err != nil {
		if errors.Is(err, storage.ErrCopyToNotEmpty) {
			return msgErrJournalCopy, nil
		}

		r.logger.Error(
			"journal copy command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return msgErrInternal, nil
	}

	return fmt.Sprintf(msgJournalCopied, cnt), nil
}

func (r *CmdProcessor) journalReportMealCommand(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid journal rm command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
	}

	// Get list from DB and user settings
	ctx, cancel := context.WithTimeout(context.Background(), _stgOperationTimeout*2)
	defer cancel()

	lst, err := r.stg.GetJournalMealReport(ctx, userID, ts, storage.NewMealFromString(cmdParts[1]))
	if err != nil {
		if errors.Is(err, storage.ErrJournalMealReportEmpty) {
			return msgErrEmptyList, nil
		}

		r.logger.Error(
			"journal rm command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return msgErrInternal, nil
	}

	var sb strings.Builder
	var totalCal float64
	for _, item := range lst {
		foodLbl := item.FoodName
		if item.FoodBrand != "" {
			foodLbl += " - " + item.FoodBrand
		}
		sb.WriteString(fmt.Sprintf("<b>%s [%s]</b>:\n", foodLbl, item.FoodKey))
		sb.WriteString(fmt.Sprintf("%.1f г., %.2f ккал\n", item.FoodWeight, item.Cal))
		totalCal += item.Cal
	}
	sb.WriteString(fmt.Sprintf("\n<b>Всего, ккал:</b> %.2f", totalCal))

	return sb.String(), []any{&tele.SendOptions{ParseMode: tele.ModeHTML}}
}

func (r *CmdProcessor) journalReportDayCommand(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid journal rd command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
	}
	tsStr := formatTimestamp(ts)

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

			return msgErrInternal, nil
		}
	}

	lst, err := r.stg.GetJournalReport(ctx, userID, ts, ts)
	if err != nil {
		if errors.Is(err, storage.ErrJournalReportEmpty) {
			return msgErrEmptyList, nil
		}

		r.logger.Error(
			"journal rd command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return msgErrInternal, nil
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
						calDiffSnippet(us, totalCal),
					),
					html.Attrs{"colspan": "6"}))).
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
	return &tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("report_%s.html", tsStr),
	}, nil
}

func (r *CmdProcessor) journalReportWeek(cmdParts []string, userID int64) (any, []any) {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid journal rw command",
			zap.String("reason", "len parts"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
		)
		return msgErrInvalidCommand, nil
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
		return msgErrInvalidCommand, nil
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

			return msgErrInternal, nil
		}
	}

	lst, err := r.stg.GetJournalStats(ctx, userID, tsStart, tsEnd)
	if err != nil {
		if errors.Is(err, storage.ErrJournalStatsEmpty) {
			return msgErrEmptyList, nil
		}

		r.logger.Error(
			"journal rw command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return msgErrInternal, nil
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

		return msgErrInternal, nil
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

	return &tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("stats_%s_%s.html", tsStartStr, tsEndStr),
	}, nil
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

func pfcSnippet(val, totalVal float64) html.IELement {
	var s string

	if totalVal == 0 {
		s = fmt.Sprintf("%.2f", val)
	} else {
		s = fmt.Sprintf("%.2f (%.2f%%)", val, val/totalVal*100)
	}

	return html.NewS(s)
}
