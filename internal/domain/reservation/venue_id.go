package reservation

import "fmt"

var (
	ErrVenueIDInvalid  = fmt.Errorf("invalid venue ID")
	ErrVenueIDZero     = fmt.Errorf("venue ID can't be zero")
	ErrVenueIDNegative = fmt.Errorf("venue ID can't be negative")
)

type VenueID int

func (id VenueID) Validate() error {
	if id == 0 {
		return fmt.Errorf("%w: %w", ErrVenueIDInvalid, ErrVenueIDZero)
	}

	if id < 0 {
		return fmt.Errorf("%w: %w", ErrVenueIDInvalid, ErrVenueIDNegative)
	}

	return nil
}
