package reservation

import (
	"fmt"
	"time"
)

var (
	ErrTimeRangeInvalid        = fmt.Errorf("invalid time range")
	ErrTimeRangeStartIsZero    = fmt.Errorf("Time range start can't be zero")
	ErrTimeRangeEndIsZero      = fmt.Errorf("Time range end can't be zero")
	ErrTimeRangeEndBeforeStart = fmt.Errorf("Time range end can't be before start")
	ErrTimeRangeSameEndStart   = fmt.Errorf("Time range start and end can't be the same time")
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
		return nil, fmt.Errorf("%w: %w", ErrTimeRangeInvalid, err)
	}

	return r, nil
}

func (r *TimeRange) Duration() time.Duration {
	return r.End.Sub(r.Start)
}

func (r *TimeRange) Validate() error {
	if r.Start.IsZero() {
		return ErrTimeRangeStartIsZero
	}

	if r.End.IsZero() {
		return ErrTimeRangeEndIsZero
	}

	if r.End.Before(r.Start) {
		return ErrTimeRangeEndBeforeStart
	}

	if r.Start.Equal(r.End) {
		return ErrTimeRangeSameEndStart
	}

	return nil
}
