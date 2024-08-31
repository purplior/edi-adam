package mydate

import (
	"strconv"
	"strings"
	"time"
)

func FormatDate(t time.Time) string {
	yStamp := strconv.Itoa(t.Year())
	mStamp := toDigit(strconv.Itoa(int(t.Month())), 2)
	dStamp := toDigit(strconv.Itoa(t.Day()), 2)

	return yStamp + "." + mStamp + "." + dStamp
}

func FormatParseDate(dateStamp string) time.Time {
	segments := strings.Split(dateStamp, "-")
	segmentLength := len(segments)

	if segmentLength != 3 {
		return Zero()
	}

	year, err := strconv.Atoi(segments[0])
	if err != nil {
		return Zero()
	}
	month, err := strconv.Atoi(segments[1])
	if err != nil {
		return Zero()
	}
	day, err := strconv.Atoi(segments[2])
	if err != nil {
		return Zero()
	}

	return DayStart(year, month, day)
}

func toDigit(str string, length int) string {
	result := str

	for len(result) < length {
		result = "0" + result
	}

	return result
}
