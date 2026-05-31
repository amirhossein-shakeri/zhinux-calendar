package reservation

import (
	"fmt"
	"time"
)

// Represents a time period with clear start and end
type TimeRange struct {
	Start time.Time
	End   time.Time
}

func NewTimeRange(start, end time.Time) (*TimeRange, error) {
	if start.IsZero() {
		return nil, fmt.Errorf("Time range start can't be zero")
	}

	if end.IsZero() {
		return nil, fmt.Errorf("Time range end can't be zero")
	}

	if end.Before(start) {
		return nil, fmt.Errorf("Time range end can't be before start")
	}

	if start.Equal(end) {
		return nil, fmt.Errorf("Time range start and end can't be the same time")
	}

	return &TimeRange{
		Start: start,
		End:   end,
	}, nil
}
