package util

import (
	"time"
)

func IsValidMonth(month time.Month, input time.Time) bool {
	first, last := MonthInterval(input.Year(), month)
	if input.After(first) && input.Before(last) {
		return true
	} else {
		return false
	}
}

func MonthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC), time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
}
