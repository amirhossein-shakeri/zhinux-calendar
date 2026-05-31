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
	r := &TimeRange{
		Start: start,
		End:   end,
	}

	if err := r.Validate(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *TimeRange) Validate() error {
	if r.Start.IsZero() {
		return fmt.Errorf("Time range start can't be zero")
	}

	if r.End.IsZero() {
		return fmt.Errorf("Time range end can't be zero")
	}

	if r.End.Before(r.Start) {
		return fmt.Errorf("Time range end can't be before start")
	}

	if r.Start.Equal(r.End) {
		return fmt.Errorf("Time range start and end can't be the same time")
	}

	return nil
}
