package reservation_test

import (
	"errors"
	"testing"

	"github.com/amirhossein-shakeri/zhinux-calendar/internal/domain/reservation"
)

func TestVenueIDValidation(t *testing.T) {
	cases := []struct {
		name string
		ID   reservation.VenueID
		want error
	}{
		{"valid positive venue ID", reservation.VenueID(153), nil},
		{"invalid zero venue ID", reservation.VenueID(0), reservation.ErrVenueIDZero},
		{"invalid negative venue ID", reservation.VenueID(-76), reservation.ErrVenueIDNegative},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(ttt *testing.T) {
			got := tt.ID.Validate()
			if tt.want != nil && !errors.Is(got, tt.want) {
				ttt.Errorf("Validate() got %v; want %v", got, tt.want)
			}
			if tt.want != nil && !errors.Is(got, reservation.ErrVenueIDInvalid) {
				ttt.Errorf("Validate() got %v; expected %v", got, reservation.ErrVenueIDInvalid)
			}
		})
	}
}
