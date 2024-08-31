package mydate

import "time"

func ConvertUnixMilliToTime(unix int) time.Time {
	return time.UnixMilli(int64(unix))
}
