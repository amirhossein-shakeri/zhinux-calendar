package reservation

type VenueID int

func (v VenueID) IsValid() bool {
	return v > 0
}
