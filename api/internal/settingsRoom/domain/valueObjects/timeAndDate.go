package valueobjects

import (
	"time"
)

type (
	TimeAndDate struct {
		Time *time.Time
		Date *time.Time
	}
)

func NewTimeAndDate(time, date *time.Time) (*TimeAndDate, error) {

	return &TimeAndDate{
		Time: time,
		Date: date,
	}, nil
}
