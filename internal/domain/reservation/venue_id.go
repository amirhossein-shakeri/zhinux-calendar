package reservation

type VenueID int
type VenuePublicID string // UUIDv7

func (v VenueID) IsValid() bool {
	return v > 0
}
