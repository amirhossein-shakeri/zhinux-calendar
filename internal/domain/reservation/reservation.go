package reservation

type Reservation struct {
	ID int

	Slot

	IdempotencyKey string
}
