package cmdproc

import (
	"time"
)

const (
	_cssBotstrapURL = "https://devldavydov.github.io/css/bootstrap/bootstrap.min.css"
	_jsBootstrapURL = "https://devldavydov.github.io/js/bootstrap/bootstrap.bundle.min.js"
	_jsChartURL     = "https://devldavydov.github.io/js/chartjs/chart.umd.min.js"
)

func parseTimestamp(sTimestamp string) (int64, error) {
	t, err := time.Parse("02.01.2006", sTimestamp)
	if err != nil {
		return 0, err
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Unix(), nil
}

func formatTimestamp(tsUnix int64) string {
	ts := time.Unix(tsUnix, 0)
	return ts.Format("02.01.2006")
}

func isStartOfWeek(tsUnix int64) bool {
	return time.Unix(tsUnix, 0).Weekday() == time.Monday
}
