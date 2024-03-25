package cmdproc

import (
	"time"
)

const (
	_styleURL = "https://devldavydov.github.io/css/bootstrap/bootstrap.min.css"
	_jsURL    = "https://devldavydov.github.io/js/bootstrap/bootstrap.bundle.min.js"
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
