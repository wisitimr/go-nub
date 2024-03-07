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

func Month(input time.Time) string {
	switch input.Month() {
	case 1:
		return "มค"
	case 2:
		return "กพ"
	case 3:
		return "มีค"
	case 4:
		return "เมย"
	case 5:
		return "พค"
	case 6:
		return "มิย"
	case 7:
		return "กค"
	case 8:
		return "สค"
	case 9:
		return "กย"
	case 10:
		return "ตค"
	case 11:
		return "พศ"
	case 12:
		return "ธค"
	}
	return ""
}

func MonthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC), time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
}
