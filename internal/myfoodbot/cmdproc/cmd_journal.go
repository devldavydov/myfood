package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/devldavydov/myfood/internal/common/html"
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
	case "cp":
		return r.journalCopyCommand(c, cmdParts[1:], userID)
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

func (r *CmdProcessor) journalCopyCommand(c tele.Context, cmdParts []string, userID int64) error {
	if len(cmdParts) != 4 {
		r.logger.Error(
			"invalid journal copy command",
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
			"invalid journal copy command",
			zap.String("reason", "ts from format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
	}

	tsTo, err := parseTimestamp(cmdParts[2])
	if err != nil {
		r.logger.Error(
			"invalid journal copy command",
			zap.String("reason", "ts to format"),
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)
		return c.Send(msgErrInvalidCommand)
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
			return c.Send(msgErrJournalCopy)
		}

		r.logger.Error(
			"journal copy command DB error",
			zap.Strings("command", cmdParts),
			zap.Int64("userid", userID),
			zap.Error(err),
		)

		return c.Send(msgErrInternal)
	}

	return c.Send(fmt.Sprintf(msgJournalCopied, cnt))
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
	return c.Send(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
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

	tsStart, err := parseTimestamp(cmdParts[0])
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

		return c.Send(msgErrInternal)
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

	return c.Send(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("stats_%s_%s.html", tsStartStr, tsEndStr),
	})
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
