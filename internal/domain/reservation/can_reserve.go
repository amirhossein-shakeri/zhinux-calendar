package reservation

func CanReserve(
	existing []Reservation,
	candidate Reservation,
) (allowed bool, reason string) {
	// Loop through existing reservations
	for _, r := range existing {
		// Skip reservations of other venue IDs
		if r.VenueID != candidate.VenueID {
			continue
		}

		// Not allowed if candidate overlaps at least one existing reservation
		if candidate.TimeRange.Overlaps(r.TimeRange) {
			if candidate.Slot.IsEqualTo(&r.Slot) && candidate.IdempotencyKey == r.IdempotencyKey {
				// Don't return duplicate error message for the same reservation/retry
				return true, ""
			} else {
				// Conflict. Candidate overlaps, the slot is not the same
				// as existing reservation, different idempotency keys
				return false, "Time range conflicts with an existing reservation on the venue"
			}
		}
	}

	// No conflicts, No retry/dup
	return true, ""
}
