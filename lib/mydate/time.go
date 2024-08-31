package mydate

import (
	"time"
)

func Zero() time.Time {
	var t time.Time

	return t
}

func Now() time.Time {
	return time.Now().In(CurrentLocation())
}

func After(d time.Duration) time.Time {
	return Now().Add(d)
}

func DayStart(year int, month int, day int) time.Time {
	t := time.Date(
		year,
		time.Month(month),
		day,
		0,
		0,
		0,
		0,
		CurrentLocation(),
	)

	return t
}

func DayEnd(year int, month int, day int) time.Time {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, CurrentLocation())
	t = t.Add((time.Duration(24) * time.Hour) - (time.Duration(1) * time.Nanosecond))

	return t
}

func DayStartByDate(date Date) time.Time {
	return DayStart(date.Year, date.Month, date.Day)
}

func DayEndByDate(date Date) time.Time {
	return DayEnd(date.Year, date.Month, date.Day)
}

func DayStartByTime(t time.Time) time.Time {
	return DayStart(t.Year(), int(t.Month()), t.Day())
}

func DayEndByTime(t time.Time) time.Time {
	return DayEnd(t.Year(), int(t.Month()), t.Day())
}

func DayStartFromNow(dayDistance int) time.Time {
	t := Now().Add(time.Duration(24*dayDistance) * time.Hour)

	return DayStart(t.Year(), int(t.Month()), t.Day())
}

func DayEndFromNow(dayDistance int) time.Time {
	t := Now().Add(time.Duration(24*dayDistance) * time.Hour)

	return DayEnd(t.Year(), int(t.Month()), t.Day())
}
