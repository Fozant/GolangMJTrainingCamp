package utils

import (
	"errors"
	"time"
)

const DateFormat = "2006-01-02"

func ToFormattedDate(t time.Time) string {
	return t.Format(DateFormat)
}

func ParseFormattedDate(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected yyyy-mm-dd")
	}
	return parsedDate, nil
}
