package verification

import "time"

type (
	QueryOption struct {
		ID   uint
		Hash string

		CreatedAtStart time.Time
		CreatedAtEnd   time.Time
	}
)
