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

func TestNewSlot_InvalidMinDuration(t *testing.T) {
	invalidVenueID := reservation.VenueID(184)
	start := time.Date(2026, 06, 04, 0, 0, 0, 0, time.UTC)
	end := start.Add(reservation.SlotMinDuration - time.Second)
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

	if !errors.Is(err, reservation.ErrSlotMinDuration) {
		t.Errorf("Got error %v; expected error %v", err, reservation.ErrSlotMinDuration)
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

func TestSlotEqual(t *testing.T) {
	start := time.Date(2026, 06, 25, 10, 54, 0, 0, time.UTC)
	rangeA, errA := reservation.NewTimeRange(start, start.Add(30*time.Minute))
	rangeB, errB := reservation.NewTimeRange(start, start.Add(60*time.Minute))
	if errA != nil || errB != nil {
		t.Fatalf("Failed to initialize time ranges: %v, %v", errA, errB)
	}

	a, errA := reservation.NewSlot(reservation.VenueID(12), rangeA)
	b, errB := reservation.NewSlot(reservation.VenueID(12), rangeB)
	c, errC := reservation.NewSlot(reservation.VenueID(14), rangeA)
	if errA != nil || errB != nil || errC != nil {
		t.Fatalf("Failed to initialize time ranges: %v, %v, %v", errA, errB, errC)
	}

	// Nil comparison
	if a.IsEqualTo(nil) {
		t.Errorf("Comparing slot with nil should be false, got true")
	}

	// Different VenueIDs, Different Ranges
	if c.IsEqualTo(b) {
		t.Errorf("Slot comparison with different VenueIDs and different TimeRanges got true, expected false")
	}

	// Same VenueIDs, Different Ranges
	if a.IsEqualTo(b) {
		t.Errorf("Slot comparison with same VenueIDs and different TimeRanges got true, expected false")
	}

	// Different VenueIDs, Same Ranges
	if a.IsEqualTo(c) {
		t.Errorf("Slot comparison with different VenueIDs and same TimeRanges got true, expected false")
	}

	// Self comparison
	if !c.IsEqualTo(c) {
		t.Errorf("Slot comparison with same VenueIDs and same TimeRanges got false, expected true")
	}
}
