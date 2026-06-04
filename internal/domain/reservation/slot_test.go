package reservation_test

import (
	"errors"
	"testing"
	"time"

	"github.com/amirhossein-shakeri/zhinux-calendar/internal/domain/reservation"
)

func TestNewSlot_InvalidTimeRange(t *testing.T) {
	validVenueID := reservation.VenueID(769)
	invalidTimeRange := reservation.TimeRange{}

	got, err := reservation.NewSlot(validVenueID, invalidTimeRange)
	if err == nil {
		t.Errorf("Got no error; expected init/validation error")
	}

	if got != nil {
		t.Errorf("Got non-nil slot; expected nil: %+v", got)
	}

	if !errors.Is(err, reservation.ErrSlotInitFailed) {
		t.Errorf("Got error %v; expected error %v", err, reservation.ErrSlotInitFailed)
	}

	if !errors.Is(err, reservation.ErrSlotInvalid) {
		t.Errorf("Got error %v; expected error %v", err, reservation.ErrSlotInvalid)
	}

	if !errors.Is(err, reservation.ErrTimeRangeInvalid) {
		t.Errorf("Got error %v; expected error %v", err, reservation.ErrTimeRangeInvalid)
	}
}

func TestNewSlot_InvalidVenueID(t *testing.T) {
	invalidVenueID := reservation.VenueID(-79)
	start := time.Date(2026, 06, 04, 0, 0, 0, 0, time.UTC)
	end := start.Add(30 * time.Minute)
	validTimeRange, err := reservation.NewTimeRange(start, end)
	if err != nil {
		t.Fatalf("Couldn't initialize valid time range: %v", err)
	}

	got, err := reservation.NewSlot(invalidVenueID, validTimeRange)
	if err == nil {
		t.Errorf("Got no error; expected init/validation error")
	}

	if got != nil {
		t.Errorf("Got non-nil slot; expected nil: %+v", got)
	}

	if !errors.Is(err, reservation.ErrSlotInitFailed) {
		t.Errorf("Got error %v; expected error %v", err, reservation.ErrSlotInitFailed)
	}

	if !errors.Is(err, reservation.ErrSlotInvalid) {
		t.Errorf("Got error %v; expected error %v", err, reservation.ErrSlotInvalid)
	}

	if !errors.Is(err, reservation.ErrVenueIDInvalid) {
		t.Errorf("Got error %v; expected error %v", err, reservation.ErrVenueIDInvalid)
	}
}

func TestNewSlot_ValidTimeRangeAndVenueID(t *testing.T) {
	validVenueID := reservation.VenueID(83)
	start := time.Date(2026, 06, 04, 0, 0, 0, 0, time.UTC)
	end := start.Add(30 * time.Minute)
	validTimeRange, err := reservation.NewTimeRange(start, end)
	if err != nil {
		t.Fatalf("Couldn't initialize valid time range: %v", err)
	}

	got, err := reservation.NewSlot(validVenueID, validTimeRange)
	if err != nil {
		t.Errorf("Got error; expected no error: %v", err)
	}

	if got == nil {
		t.Errorf("Got nil slot; expected value")
	}
}
