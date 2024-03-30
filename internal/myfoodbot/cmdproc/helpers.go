package cmdproc

import (
	"time"
)

const (
	_cssBotstrapURL = "https://devldavydov.github.io/css/bootstrap/bootstrap.min.css"
	_jsBootstrapURL = "https://devldavydov.github.io/js/bootstrap/bootstrap.bundle.min.js"
	_jsChartURL     = "https://devldavydov.github.io/js/chartjs/chart.umd.min.js"
)

func parseTimestampAsUnix(sTimestamp string) (int64, error) {
	t, err := parseTimestamp(sTimestamp)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func formatTimestampUnix(tsUnix int64) string {
	ts := time.Unix(tsUnix, 0)
	return ts.Format("02.01.2006")
}

func parseTimestamp(sTimestamp string) (time.Time, error) {
	t, err := time.Parse("02.01.2006", sTimestamp)
	if err != nil {
		return time.Time{}, err
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
