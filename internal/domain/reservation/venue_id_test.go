package reservation_test

import (
	"testing"

	"github.com/amirhossein-shakeri/zhinux-calendar/internal/domain/reservation"
)

func TestVenueIDValidation(t *testing.T) {
	cases := []struct {
		name string
		ID   reservation.VenueID
		want bool
	}{
		{"valid positive venue ID", reservation.VenueID(153), true},
		{"invalid zero venue ID", reservation.VenueID(0), false},
		{"invalid negative venue ID", reservation.VenueID(-76), false},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(ttt *testing.T) {
			got := tt.ID.IsValid()
			if got != tt.want {
				ttt.Errorf("IsValid() got %v; want %v", got, tt.want)
			}
		})
	}
}
