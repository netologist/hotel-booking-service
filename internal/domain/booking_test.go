package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/domain"
)

const (
	testJaneDoe   = "Jane Doe"
	testJaneEmail = "jane@example.com"
)

func TestNewBooking_Success(t *testing.T) {
	t.Parallel()
	hotelID := uuid.New()
	room, err := domain.NewRoom(hotelID, "201", domain.RoomTypeDouble, 2)
	require.NoError(t, err)

	checkIn := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	checkOut := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)

	b, err := domain.NewBooking(room, testJaneDoe, testJaneEmail, checkIn, checkOut, 2)
	require.NoError(t, err)

	assert.NotEqual(t, uuid.Nil, b.ID)
	assert.NotEmpty(t, b.Reference)
	assert.Equal(t, room.ID, b.RoomID)
	assert.Equal(t, hotelID, b.HotelID)
	assert.Equal(t, testJaneDoe, b.GuestName)
	assert.Equal(t, testJaneEmail, b.GuestEmail)
	assert.Equal(t, 2, b.NumberOfGuests)
	assert.True(t, b.CheckInDate.Equal(checkIn))
	assert.True(t, b.CheckOutDate.Equal(checkOut))
}

func TestNewBooking_Invariants(t *testing.T) { //nolint:funlen // table covers all domain invariant variations
	t.Parallel()
	hotelID := uuid.New()
	validRoom, err := domain.NewRoom(hotelID, "101", domain.RoomTypeSingle, 1)
	require.NoError(t, err)

	checkIn := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	checkOut := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		room       domain.Room
		guestName  string
		guestEmail string
		checkIn    time.Time
		checkOut   time.Time
		guests     int
		wantErr    error
	}{
		{
			name:       "nil room id",
			room:       domain.Room{},
			guestName:  testJaneDoe,
			guestEmail: testJaneEmail,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     1,
			wantErr:    domain.ErrInvalidBooking,
		},
		{
			name:       "short guest name",
			room:       validRoom,
			guestName:  "J",
			guestEmail: testJaneEmail,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     1,
			wantErr:    domain.ErrInvalidBooking,
		},
		{
			name:       "invalid email",
			room:       validRoom,
			guestName:  testJaneDoe,
			guestEmail: "notanemail",
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     1,
			wantErr:    domain.ErrInvalidBooking,
		},
		{
			name:       "check-out before check-in",
			room:       validRoom,
			guestName:  testJaneDoe,
			guestEmail: testJaneEmail,
			checkIn:    checkOut,
			checkOut:   checkIn,
			guests:     1,
			wantErr:    domain.ErrInvalidBooking,
		},
		{
			name:       "check-out equals check-in",
			room:       validRoom,
			guestName:  testJaneDoe,
			guestEmail: testJaneEmail,
			checkIn:    checkIn,
			checkOut:   checkIn,
			guests:     1,
			wantErr:    domain.ErrInvalidBooking,
		},
		{
			name:       "zero guests",
			room:       validRoom,
			guestName:  testJaneDoe,
			guestEmail: testJaneEmail,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     0,
			wantErr:    domain.ErrInvalidBooking,
		},
		{
			name:       "guests exceed room capacity",
			room:       validRoom,
			guestName:  testJaneDoe,
			guestEmail: testJaneEmail,
			checkIn:    checkIn,
			checkOut:   checkOut,
			guests:     2,
			wantErr:    domain.ErrCapacityExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewBooking(tt.room, tt.guestName, tt.guestEmail, tt.checkIn, tt.checkOut, tt.guests)
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGenerateBookingReference(t *testing.T) {
	t.Parallel()
	ref1 := domain.GenerateBookingReference()
	ref2 := domain.GenerateBookingReference()

	assert.NotEmpty(t, ref1)
	assert.NotEmpty(t, ref2)
	assert.NotEqual(t, ref1, ref2, "generated booking references must be unique")
}
