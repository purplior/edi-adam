package mydate

import "time"

var current *time.Location

const (
	Location_Seoul = "Asia/Seoul"
)

func initLocation(location *string) error {
	var locLabel string

	if location == nil {
		locLabel = Location_Seoul
	} else {
		locLabel = *location
	}

	loc, err := time.LoadLocation(locLabel)
	if err != nil {
		return err
	}

	current = loc

	return nil
}

func CurrentLocation() *time.Location {
	if current == nil {
		if err := Init(nil); err != nil {
			panic(err)
		}
	}

	return current
}
