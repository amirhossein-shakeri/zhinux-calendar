package reservation

import (
	"fmt"
	"time"
)

var (
	ErrTimeRangeInitFailed     = fmt.Errorf("failed to initialize new time range")
	ErrTimeRangeInvalid        = fmt.Errorf("invalid time range")
	ErrTimeRangeStartIsZero    = fmt.Errorf("Time range start can't be zero")
	ErrTimeRangeEndIsZero      = fmt.Errorf("Time range end can't be zero")
	ErrTimeRangeEndBeforeStart = fmt.Errorf("Time range end can't be before start")
	ErrTimeRangeSameEndStart   = fmt.Errorf("Time range start and end can't be the same time")
)

// Represents a time period with clear start and end
type TimeRange struct {
	start time.Time // Immutable value object member
	end   time.Time // Immutable value object member
}

// NewTimeRange constructs a new TimeRange with invariants enforced
// and since it's small, we'd prefer returning value over pointer
func NewTimeRange(start, end time.Time) (TimeRange, error) {
	r := TimeRange{
		start: start,
		end:   end,
	}

	if err := r.Validate(); err != nil {
		return TimeRange{}, fmt.Errorf("%w: %w", ErrTimeRangeInitFailed, err)
	}

	return r, nil
}

func (r *TimeRange) Start() time.Time { return r.start }
func (r *TimeRange) End() time.Time   { return r.end }

func (r *TimeRange) Duration() time.Duration {
	return r.End().Sub(r.Start())
}

func (r *TimeRange) Overlaps(otherTimeRange TimeRange) bool {
	return r.start.Before(otherTimeRange.end) && otherTimeRange.start.Before(r.end)
}

func (r *TimeRange) Validate() error {
	if r.Start().IsZero() {
		return fmt.Errorf("%w: %w", ErrTimeRangeInvalid, ErrTimeRangeStartIsZero)
	}

	if r.End().IsZero() {
		return fmt.Errorf("%w: %w", ErrTimeRangeInvalid, ErrTimeRangeEndIsZero)
	}

	if r.End().Before(r.Start()) {
		return fmt.Errorf("%w: %w", ErrTimeRangeInvalid, ErrTimeRangeEndBeforeStart)
	}

	if r.Start().Equal(r.End()) {
		return fmt.Errorf("%w: %w", ErrTimeRangeInvalid, ErrTimeRangeSameEndStart)
	}

	return nil
}

// TODO: Add unit tests
func (r *TimeRange) IsEqualTo(otherTimeRange *TimeRange) bool {
	if otherTimeRange == nil {
		return false
	}

	return r.Start().Equal(otherTimeRange.Start()) &&
		r.End().Equal(otherTimeRange.End())
}
