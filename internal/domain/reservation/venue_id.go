package reservation

type VenueID int

func (id VenueID) IsValid() bool {
	return id > 0
}
