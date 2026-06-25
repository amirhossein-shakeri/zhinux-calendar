package reservation_test

import (
	"errors"
	"testing"
	"time"

	"github.com/amirhossein-shakeri/zhinux-calendar/internal/domain/reservation"
)

func TestTimeRangeValidation(t *testing.T) {
	var zeroTime time.Time
	validDate := time.Date(2026, 06, 03, 0, 0, 0, 0, time.UTC)
	emptyTimeRange := reservation.TimeRange{}

	cases := []struct {
		name  string
		start time.Time
		end   time.Time

		wantErr      error
		wantDuration time.Duration
	}{
		{
			name:    "zero start",
			start:   zeroTime,
			end:     validDate,
			wantErr: reservation.ErrTimeRangeStartIsZero,
		},
		{
			name:    "zero end",
			start:   validDate,
			end:     zeroTime,
			wantErr: reservation.ErrTimeRangeEndIsZero,
		},
		{
			name:    "end before start",
			start:   validDate,
			end:     validDate.Add(-1 * time.Hour),
			wantErr: reservation.ErrTimeRangeEndBeforeStart,
		},
		{
			name:    "end the same as start",
			start:   validDate,
			end:     validDate,
			wantErr: reservation.ErrTimeRangeSameEndStart,
		},
		{
			name:         "valid start and end; Ninety minutes",
			start:        validDate,
			end:          validDate.Add(90 * time.Minute),
			wantDuration: 90 * time.Minute,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(ttt *testing.T) {
			got, err := reservation.NewTimeRange(tt.start, tt.end)

			if tt.wantErr != nil {
				if err == nil {
					ttt.Fatalf("NewTimeRange(%v, %v): Got no error; want %v",
						tt.start, tt.end, tt.wantErr)
				}

				if !errors.Is(err, tt.wantErr) {
					ttt.Fatalf("NewTimeRange(%v, %v): Got %v; want %v",
						tt.start, tt.end, err, tt.wantErr)
				}

				if !errors.Is(err, reservation.ErrTimeRangeInvalid) {
					ttt.Fatalf("NewTimeRange(%v, %v): Got %v; want %v",
						tt.start, tt.end, err, reservation.ErrTimeRangeInvalid)
				}

				if !errors.Is(err, reservation.ErrTimeRangeInitFailed) {
					ttt.Fatalf("NewTimeRange(%v, %v): Got %v; want %v",
						tt.start, tt.end, err, reservation.ErrTimeRangeInitFailed)
				}

				if got != emptyTimeRange {
					ttt.Fatalf("Expected empty TimeRage, got %+v", got)
				}

				return
			}

			if err != nil {
				ttt.Fatalf("NewTimeRange(%v, %v): Unexpected error %v",
					tt.start, tt.end, err)
			}

			if got == emptyTimeRange {
				ttt.Fatalf(
					"NewTimeRange(%v, %v): Expected non-empty TimeRange, got nil",
					tt.start, tt.end)
			}

			if !got.Start().Equal(tt.start) {
				ttt.Errorf(
					"NewTimeRange(%v, %v): Mismatching starts, got %v; expected %v",
					tt.start, tt.end, got.Start(), tt.start)
			}

			if !got.End().Equal(tt.end) {
				ttt.Errorf(
					"NewTimeRange(%v, %v): Mismatching ends, got %v; expected %v",
					tt.start, tt.end, got.End(), tt.end)
			}

			duration := got.Duration()
			if tt.wantDuration != 0 && duration != tt.wantDuration {
				ttt.Errorf(
					"NewTimeRange(%v, %v): Wrong duration, got %v; want %v",
					tt.start, tt.end, duration, tt.wantDuration)
			}
		})
	}
}

// In inverval math, there are exactly 6 ways two intervals can
// relate. If we cover those 6, we've tested the entire logic
// space, so we won't need to hardcode a full matrix test case
// and a static truth table.
//
// Separated(A before B)
// Touching(A end = B start)
// Overlapping(Partial)
// Contained(A inside B)
// Containing(B inside A)
// Identical
func TestTimeRangeOverlap(t *testing.T) {
	baseStart := time.Date(2026, 06, 24, 10, 0, 0, 0, time.UTC)
	baseEnd := baseStart.Add(2 * time.Hour)

	// Reference range A: [10:00, 12:00)
	a, err := reservation.NewTimeRange(baseStart, baseEnd)
	if err != nil {
		t.Fatalf("Failed to initialize reference time range: %v", err)
	}

	cases := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected bool
	}{
		{"Separated Before", baseStart.Add(-4 * time.Hour), baseStart.Add(-2 * time.Hour), false},
		{"Touching Before", baseStart.Add(-2 * time.Hour), baseStart, false},
		{"Overlapping Start", baseStart.Add(-1 * time.Hour), baseStart.Add(1 * time.Hour), true},
		{"Identical", baseStart, baseEnd, true},
		{"Inside", baseStart.Add(30 * time.Minute), baseEnd.Add(-30 * time.Minute), true},
		{"Encompassing", baseStart.Add(-1 * time.Hour), baseEnd.Add(1 * time.Hour), true},
		{"Overlapping End", baseStart.Add(1 * time.Hour), baseEnd.Add(1 * time.Hour), true},
		{"Touching After", baseEnd, baseEnd.Add(2 * time.Hour), false},
		{"Separated After", baseEnd.Add(1 * time.Hour), baseEnd.Add(3 * time.Hour), false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b, err := reservation.NewTimeRange(tt.start, tt.end)
			if err != nil {
				t.Fatalf("Failed to initialize sub time range: %v", err)
			}

			got := a.Overlaps(b)
			if got != tt.expected {
				t.Errorf("%s: Overlaps() = %v, want %v", tt.name, got, tt.expected)
			}

			// Check symmetry invariant(Reverse should result the same)
			if a.Overlaps(b) != b.Overlaps(a) {
				t.Errorf("%s: Symmetry failed! A.Overlaps(B) != B.Overlaps(A)", tt.name)
			}
		})
	}
}

func TestTimeRangeEqual(t *testing.T) {
	start := time.Date(2026, 06, 25, 10, 43, 0, 0, time.UTC)
	a, errA := reservation.NewTimeRange(start, start.Add(30*time.Minute))
	b, errB := reservation.NewTimeRange(start, start.Add(30*time.Minute))
	c, errC := reservation.NewTimeRange(start, start.Add(60*time.Minute))
	if errA != nil || errB != nil || errC != nil {
		t.Fatalf("Failed to initialize time ranges: %v, %v, %v", errA, errB, errC)
	}

	if !a.IsEqualTo(&b) {
		t.Errorf("Time range equality check, got false; expected true")
	}

	if b.IsEqualTo(&c) {
		t.Errorf("Time range equality check, got true; expected false")
	}
}
