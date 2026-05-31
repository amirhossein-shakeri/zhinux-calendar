package reservation

import "fmt"

type Slot struct {
	VenueID   VenueID
	TimeRange TimeRange
}

func (s *Slot) Validate() error {
	if !s.VenueID.IsValid() {
		return fmt.Errorf("Invalid slot: Invalid venue ID")
	}

	if err := s.TimeRange.Validate(); err != nil {
		return fmt.Errorf("Invalid slot time range: %w", err)
	}

	return nil
}
