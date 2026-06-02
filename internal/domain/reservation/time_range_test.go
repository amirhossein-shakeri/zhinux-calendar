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

				if got != nil {
					ttt.Fatalf("Expected nil TimeRage, got %+v", got)
				}

				return
			}

			if err != nil {
				ttt.Fatalf("NewTimeRange(%v, %v): Unexpected error %v",
					tt.start, tt.end, err)
			}

			if got == nil {
				ttt.Fatalf(
					"NewTimeRange(%v, %v): Expected non-nil TimeRange, got nil",
					tt.start, tt.end)
			}

			if !got.Start.Equal(tt.start) {
				ttt.Errorf(
					"NewTimeRange(%v, %v): Mismatching starts, got %v; expected %v",
					tt.start, tt.end, got.Start, tt.start)
			}

			if !got.End.Equal(tt.end) {
				ttt.Errorf(
					"NewTimeRange(%v, %v): Mismatching ends, got %v; expected %v",
					tt.start, tt.end, got.End, tt.end)
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
