package cmdproc

import (
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	_cssBotstrapURL = "https://devldavydov.github.io/css/bootstrap/bootstrap.min.css"
	_jsBootstrapURL = "https://devldavydov.github.io/js/bootstrap/bootstrap.bundle.min.js"
	_jsChartURL     = "https://devldavydov.github.io/js/chartjs/chart.umd.min.js"
)

var optsHTML = &tele.SendOptions{ParseMode: tele.ModeHTML}

func (r *CmdProcessor) parseTimestamp(sTimestamp string) (time.Time, error) {
	var t time.Time
	var err error

	if sTimestamp == "" {
		t = time.Now().In(r.tz)
	} else {
		t, err = time.Parse("02.01.2006", sTimestamp)
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
}

func formatTimestamp(ts time.Time) string {
	return ts.Format("02.01.2006")
}

func getStartOfWeek(ts time.Time) time.Time {
	day := 24 * time.Hour

	switch ts.Weekday() {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday:
		return ts.Add(-1 * time.Duration((ts.Weekday() - 1)) * day)
	case time.Sunday:
		return ts.Add(-6 * day)
	}
	return ts
}
