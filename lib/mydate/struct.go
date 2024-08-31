package mydate

import "time"

type (
	Date struct {
		Year  int
		Month int
		Day   int
	}
)

func (d *Date) ToDayStart() time.Time {
	return DayStart(d.Year, d.Month, d.Day)
}

func (d *Date) ToDayEnd() time.Time {
	return DayEnd(d.Year, d.Month, d.Day)
}

func (d *Date) CheckSameDay(other Date) bool {
	return d.Year == other.Year && d.Month == other.Month && d.Day == other.Day
}
