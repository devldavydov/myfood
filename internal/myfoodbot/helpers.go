package myfoodbot

import (
	"time"
)

func parseTimestamp(sTimestamp string) (int64, error) {
	if sTimestamp == "" {
		t := time.Now()
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix(), nil
	} else {
		t, err := time.Parse("02.01.2006", sTimestamp)
		if err != nil {
			return 0, err
		}
		return t.Unix(), nil
	}
}
