package reservation_test

import (
	"strings"
	"testing"
	"time"

	"github.com/amirhossein-shakeri/zhinux-calendar/internal/domain/reservation"
)

func mustNewTimeRange(t *testing.T, start, end time.Time) reservation.TimeRange {
	t.Helper()

	tr, err := reservation.NewTimeRange(start, end)
	if err != nil {
		t.Fatalf("failed to create time range (%s, %s): %v", start, end, err)
	}

	return tr
}

func mustNewSlot(t *testing.T, venueID reservation.VenueID, start, end time.Time) reservation.Slot {
	t.Helper()

	slot, err := reservation.NewSlot(venueID, mustNewTimeRange(t, start, end))
	if err != nil {
		t.Fatalf("failed to create slot for venue %v (%s, %s): %v", venueID, start, end, err)
	}

	return *slot
}

func mustNewReservation(t *testing.T, id int, venueID reservation.VenueID, start, end time.Time, idempotencyKey string) reservation.Reservation {
	t.Helper()
	return reservation.Reservation{
		ID:             id,
		Slot:           mustNewSlot(t, venueID, start, end),
		IdempotencyKey: idempotencyKey,
	}
}

func assertContainsAll(t *testing.T, got string, needles []string) {
	t.Helper()

	gotLower := strings.ToLower(got)
	for _, needle := range needles {
		if !strings.Contains(gotLower, strings.ToLower(needle)) {
			t.Errorf("reason mismatch: got %q, expected to contain %q", got, needle)
		}
	}
}

func TestCanReserve(t *testing.T) {
	base := time.Date(2026, time.June, 5, 9, 0, 0, 0, time.UTC)

	t.Run("table-driven", func(t *testing.T) {
		cases := []struct {
			name             string
			existing         []reservation.Reservation
			candidate        reservation.Reservation
			wantAllowed      bool
			wantReasonNeedle []string
		}{
			{
				name:        "empty existing reservations -> allowed",
				existing:    nil,
				candidate:   mustNewReservation(t, 1, reservation.VenueID(1), base, base.Add(30*time.Minute), "k-new-1"),
				wantAllowed: true,
			},
			{
				name:        "same venue before existing window",
				existing:    []reservation.Reservation{mustNewReservation(t, 1, reservation.VenueID(1), base.Add(2*time.Hour), base.Add(3*time.Hour), "k-1")},
				candidate:   mustNewReservation(t, 2, reservation.VenueID(1), base, base.Add(90*time.Minute), "k-2"),
				wantAllowed: true,
			},
			{
				name:        "same venue touching previous end exactly",
				existing:    []reservation.Reservation{mustNewReservation(t, 1, reservation.VenueID(1), base, base.Add(30*time.Minute), "k-3")},
				candidate:   mustNewReservation(t, 2, reservation.VenueID(1), base.Add(30*time.Minute), base.Add(60*time.Minute), "k-4"),
				wantAllowed: true,
			},
			{
				name:        "same venue touching previous start exactly",
				existing:    []reservation.Reservation{mustNewReservation(t, 1, reservation.VenueID(1), base.Add(60*time.Minute), base.Add(120*time.Minute), "k-5")},
				candidate:   mustNewReservation(t, 2, reservation.VenueID(1), base, base.Add(60*time.Minute), "k-6"),
				wantAllowed: true,
			},
			{
				name:        "same venue overlapping at start",
				existing:    []reservation.Reservation{mustNewReservation(t, 1, reservation.VenueID(1), base, base.Add(2*time.Hour), "k-7")},
				candidate:   mustNewReservation(t, 2, reservation.VenueID(1), base.Add(-30*time.Minute), base.Add(30*time.Minute), "k-8"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"overlap",
				},
			},
			{
				name:        "same venue overlapping at end",
				existing:    []reservation.Reservation{mustNewReservation(t, 1, reservation.VenueID(1), base, base.Add(2*time.Hour), "k-9")},
				candidate:   mustNewReservation(t, 2, reservation.VenueID(1), base.Add(90*time.Minute), base.Add(150*time.Minute), "k-10"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"overlap",
				},
			},
			{
				name:        "same venue fully contained by existing",
				existing:    []reservation.Reservation{mustNewReservation(t, 1, reservation.VenueID(1), base, base.Add(3*time.Hour), "k-11")},
				candidate:   mustNewReservation(t, 2, reservation.VenueID(1), base.Add(30*time.Minute), base.Add(60*time.Minute), "k-12"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"overlap",
				},
			},
			{
				name:        "same venue containing existing",
				existing:    []reservation.Reservation{mustNewReservation(t, 1, reservation.VenueID(1), base.Add(30*time.Minute), base.Add(60*time.Minute), "k-13")},
				candidate:   mustNewReservation(t, 2, reservation.VenueID(1), base, base.Add(2*time.Hour), "k-14"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"overlap",
				},
			},
			{
				name:        "same venue exact same interval",
				existing:    []reservation.Reservation{mustNewReservation(t, 1, reservation.VenueID(1), base.Add(15*time.Minute), base.Add(90*time.Minute), "k-15")},
				candidate:   mustNewReservation(t, 2, reservation.VenueID(1), base.Add(15*time.Minute), base.Add(90*time.Minute), "k-16"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"overlap",
				},
			},
			{
				name:        "different venue no conflict even with overlap",
				existing:    []reservation.Reservation{mustNewReservation(t, 1, reservation.VenueID(1), base, base.Add(2*time.Hour), "k-17")},
				candidate:   mustNewReservation(t, 2, reservation.VenueID(2), base.Add(30*time.Minute), base.Add(60*time.Minute), "k-18"),
				wantAllowed: true,
			},
			{
				name: "multiple venues mixed + same-venue conflict",
				existing: []reservation.Reservation{
					mustNewReservation(t, 1, reservation.VenueID(10), base, base.Add(45*time.Minute), "k-19"),
					mustNewReservation(t, 2, reservation.VenueID(11), base.Add(30*time.Minute), base.Add(90*time.Minute), "k-20"),
					mustNewReservation(t, 3, reservation.VenueID(12), base.Add(75*time.Minute), base.Add(120*time.Minute), "k-21"),
				},
				candidate:   mustNewReservation(t, 4, reservation.VenueID(11), base.Add(20*time.Minute), base.Add(40*time.Minute), "k-22"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"overlap",
				},
			},
			{
				name: "multiple venues mixed + no same-venue conflict",
				existing: []reservation.Reservation{
					mustNewReservation(t, 1, reservation.VenueID(10), base, base.Add(45*time.Minute), "k-23"),
					mustNewReservation(t, 2, reservation.VenueID(11), base.Add(30*time.Minute), base.Add(90*time.Minute), "k-24"),
					mustNewReservation(t, 3, reservation.VenueID(12), base.Add(75*time.Minute), base.Add(120*time.Minute), "k-25"),
				},
				candidate:   mustNewReservation(t, 4, reservation.VenueID(10), base.Add(50*time.Minute), base.Add(75*time.Minute), "k-26"),
				wantAllowed: true,
			},
			{
				name:        "exact minimum duration slot accepted",
				existing:    nil,
				candidate:   mustNewReservation(t, 1, reservation.VenueID(1), base, base.Add(reservation.SlotMinDuration), "k-27"),
				wantAllowed: true,
			},
			{
				name:        "slot below minimum duration rejected",
				existing:    nil,
				candidate:   mustNewReservation(t, 2, reservation.VenueID(1), base, base.Add(reservation.SlotMinDuration-time.Second), "k-28"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"minimum",
				},
			},
			{
				name:     "candidate with invalid venue-like payload (empty range rejected)",
				existing: nil,
				candidate: reservation.Reservation{
					ID: 3,
					Slot: reservation.Slot{
						VenueID:   reservation.VenueID(0),
						TimeRange: reservation.TimeRange{},
					},
					IdempotencyKey: "k-29",
				},
				wantAllowed: false,
				wantReasonNeedle: []string{
					"venue",
					"invalid",
				},
			},
			{
				name:     "candidate with invalid times rejected",
				existing: nil,
				candidate: reservation.Reservation{
					ID: 4,
					Slot: reservation.Slot{
						VenueID:   reservation.VenueID(1),
						TimeRange: reservation.TimeRange{},
					},
					IdempotencyKey: "k-30",
				},
				wantAllowed: false,
				wantReasonNeedle: []string{
					"invalid",
				},
			},
			{
				name: "same idempotency key + identical payload in same venue should be idempotent",
				existing: []reservation.Reservation{
					mustNewReservation(t, 100, reservation.VenueID(4), base.Add(30*time.Minute), base.Add(90*time.Minute), "idempotent-key-1"),
				},
				candidate:   mustNewReservation(t, 101, reservation.VenueID(4), base.Add(30*time.Minute), base.Add(90*time.Minute), "idempotent-key-1"),
				wantAllowed: true,
			},
			{
				name: "same idempotency key + same payload but different reservation ID is still duplicate replay",
				existing: []reservation.Reservation{
					mustNewReservation(t, 1_000_001, reservation.VenueID(5), base.Add(120*time.Minute), base.Add(180*time.Minute), "idempotent-key-2"),
				},
				candidate:   mustNewReservation(t, 2_000_002, reservation.VenueID(5), base.Add(120*time.Minute), base.Add(180*time.Minute), "idempotent-key-2"),
				wantAllowed: true,
			},
			{
				name: "same idempotency key + different time should fail fast as mismatch",
				existing: []reservation.Reservation{
					mustNewReservation(t, 20, reservation.VenueID(6), base, base.Add(45*time.Minute), "idempotent-key-3"),
				},
				candidate:   mustNewReservation(t, 21, reservation.VenueID(6), base.Add(15*time.Minute), base.Add(60*time.Minute), "idempotent-key-3"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"idempotency",
					"mismatch",
				},
			},
			{
				name: "same idempotency key + different venue should be mismatch",
				existing: []reservation.Reservation{
					mustNewReservation(t, 22, reservation.VenueID(6), base, base.Add(45*time.Minute), "idempotent-key-4"),
				},
				candidate:   mustNewReservation(t, 23, reservation.VenueID(7), base, base.Add(45*time.Minute), "idempotent-key-4"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"idempotency",
					"mismatch",
				},
			},
			{
				name: "no overlap when candidate nested outside all but existing invalid entry should fail only if implementation validates inputs",
				existing: []reservation.Reservation{
					{
						ID: 99,
						Slot: reservation.Slot{
							VenueID:   reservation.VenueID(8),
							TimeRange: reservation.TimeRange{},
						},
						IdempotencyKey: "bad-existing-1",
					},
					mustNewReservation(t, 24, reservation.VenueID(8), base.Add(300*time.Minute), base.Add(330*time.Minute), "k-39"),
				},
				candidate:   mustNewReservation(t, 25, reservation.VenueID(8), base.Add(10*time.Minute), base.Add(40*time.Minute), "k-40"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"invalid",
				},
			},
			{
				name: "unsorted existing reservations still evaluated correctly",
				existing: []reservation.Reservation{
					mustNewReservation(t, 26, reservation.VenueID(9), base.Add(4*time.Hour), base.Add(5*time.Hour), "k-41"),
					mustNewReservation(t, 27, reservation.VenueID(9), base.Add(2*time.Hour), base.Add(3*time.Hour), "k-42"),
					mustNewReservation(t, 28, reservation.VenueID(9), base.Add(1*time.Hour), base.Add(1*time.Hour+30*time.Minute), "k-43"),
				},
				candidate:   mustNewReservation(t, 29, reservation.VenueID(9), base.Add(2*time.Hour+15*time.Minute), base.Add(2*time.Hour+45*time.Minute), "k-44"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"overlap",
				},
			},
			{
				name: "boundary at candidate end equals existing start",
				existing: []reservation.Reservation{
					mustNewReservation(t, 30, reservation.VenueID(10), base.Add(2*time.Hour), base.Add(3*time.Hour), "k-45"),
				},
				candidate:   mustNewReservation(t, 31, reservation.VenueID(10), base.Add(time.Hour), base.Add(2*time.Hour), "k-46"),
				wantAllowed: true,
			},
			{
				name: "boundary at candidate start equals existing end",
				existing: []reservation.Reservation{
					mustNewReservation(t, 32, reservation.VenueID(10), base.Add(time.Hour), base.Add(2*time.Hour), "k-47"),
				},
				candidate:   mustNewReservation(t, 33, reservation.VenueID(10), base.Add(2*time.Hour), base.Add(3*time.Hour), "k-48"),
				wantAllowed: true,
			},
			{
				name: "existing and candidate overlap at far-left with different idempotency keys",
				existing: []reservation.Reservation{
					mustNewReservation(t, 34, reservation.VenueID(11), base.Add(30*time.Minute), base.Add(90*time.Minute), "k-49"),
				},
				candidate:   mustNewReservation(t, 35, reservation.VenueID(11), base.Add(20*time.Minute), base.Add(40*time.Minute), "k-50"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"overlap",
				},
			},
			{
				name: "same venue non-overlapping around a middle gap",
				existing: []reservation.Reservation{
					mustNewReservation(t, 36, reservation.VenueID(12), base, base.Add(20*time.Minute), "k-51"),
					mustNewReservation(t, 37, reservation.VenueID(12), base.Add(40*time.Minute), base.Add(60*time.Minute), "k-52"),
				},
				candidate:   mustNewReservation(t, 38, reservation.VenueID(12), base.Add(20*time.Minute), base.Add(40*time.Minute), "k-53"),
				wantAllowed: true,
			},
			{
				name: "candidate overlaps multiple same-venue reservations",
				existing: []reservation.Reservation{
					mustNewReservation(t, 39, reservation.VenueID(13), base, base.Add(30*time.Minute), "k-54"),
					mustNewReservation(t, 40, reservation.VenueID(13), base.Add(45*time.Minute), base.Add(75*time.Minute), "k-55"),
					mustNewReservation(t, 41, reservation.VenueID(13), base.Add(90*time.Minute), base.Add(120*time.Minute), "k-56"),
				},
				candidate:   mustNewReservation(t, 42, reservation.VenueID(13), base.Add(25*time.Minute), base.Add(55*time.Minute), "k-57"),
				wantAllowed: false,
				wantReasonNeedle: []string{
					"overlap",
				},
			},
		}

		if len(cases) < 25 {
			t.Fatalf("sanity check: expected at least 25 table-driven cases, got %d", len(cases))
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				gotAllowed, gotReason := reservation.CanReserve(tt.existing, tt.candidate)

				if gotAllowed != tt.wantAllowed {
					t.Fatalf("CanReserve() got (%v, %q), want allowed=%v", gotAllowed, gotReason, tt.wantAllowed)
				}

				if !tt.wantAllowed && len(tt.wantReasonNeedle) == 0 {
					t.Errorf("expecting a rejected reason but got empty wantReasonNeedles")
				}

				if len(tt.wantReasonNeedle) > 0 {
					if gotReason == "" {
						t.Fatalf("CanReserve() got empty reason for denied request: expected %v", tt.wantReasonNeedle)
					}
					assertContainsAll(t, gotReason, tt.wantReasonNeedle)
				}

				if tt.wantAllowed && len(tt.wantReasonNeedle) > 0 {
					assertContainsAll(t, gotReason, tt.wantReasonNeedle)
				}
			})
		}
	})
}
