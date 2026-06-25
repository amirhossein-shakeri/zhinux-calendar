package reservation

import (
	"fmt"
	"time"
)

const (
	SlotMinDuration time.Duration = time.Minute * 15
)

var (
	ErrSlotInitFailed  = fmt.Errorf("failed to initialize new slot")
	ErrSlotInvalid     = fmt.Errorf("invalid slot")
	ErrSlotMinDuration = fmt.Errorf("slot can't be shorter than %s", SlotMinDuration)
)

type Slot struct {
	VenueID VenueID
	TimeRange
}

func NewSlot(vID VenueID, tr TimeRange) (*Slot, error) {
	slot := &Slot{
		VenueID:   vID,
		TimeRange: tr,
	}

	if err := slot.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSlotInitFailed, err)
	}

	return slot, nil
}

func (s *Slot) Validate() error {
	if err := s.VenueID.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrSlotInvalid, err)
	}

	if err := s.TimeRange.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrSlotInvalid, err)
	}

	return nil
}

// TODO: Add unit tests
func (s *Slot) IsEqualTo(otherSlot *Slot) bool {
	if otherSlot == nil {
		return false
	}

	return s.VenueID == otherSlot.VenueID &&
		s.TimeRange.IsEqualTo(&otherSlot.TimeRange)
}
