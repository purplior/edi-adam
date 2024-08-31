package mydate

import (
	"time"
)

func DayRangeOfMonth(
	year int,
	month int,
) (firstDay int, latsDay int) {
	firstDayTime := DayStart(year, month, 1)
	lastDayTime := firstDayTime.AddDate(0, 1, -1)

	return 1, lastDayTime.Day()
}

func Between(
	from time.Time,
	to time.Time,
) []Date {
	fromStart := DayStartByTime(from)
	toStart := DayStartByTime(to)
	dayCount := int(toStart.Sub(fromStart)/(24*time.Hour)) + 1

	if dayCount <= 0 {
		return []Date{}
	}

	dates := make([]Date, dayCount)
	cursor := fromStart
	for i := 0; i < dayCount; i++ {
		dates[i] = Date{
			Year:  cursor.Year(),
			Month: int(cursor.Month()),
			Day:   cursor.Day(),
		}
		cursor = cursor.Add(time.Hour * 24)
	}

	return dates
}
