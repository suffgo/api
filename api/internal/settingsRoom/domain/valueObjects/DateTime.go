package valueobjects

import (
	"time"
)

type (
	DateTime struct {
		DateTime *time.Time
	}
)

func NewDateTime(dateTime *time.Time) (*DateTime, error) {

	return &DateTime{
		DateTime: dateTime,
	}, nil
}
