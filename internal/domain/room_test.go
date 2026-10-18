package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/domain"
)

func TestNewRoom_Success(t *testing.T) {
	t.Parallel()
	hotelID := uuid.New()

	room, err := domain.NewRoom(hotelID, "101", domain.RoomTypeSingle, 1)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, room.ID)
	assert.Equal(t, hotelID, room.HotelID)
	assert.Equal(t, "101", room.RoomNumber)
	assert.Equal(t, domain.RoomTypeSingle, room.RoomType)
	assert.Equal(t, 1, room.Capacity)
}

func TestNewRoom_Invariants(t *testing.T) { //nolint:funlen // test table covers all domain invariant variations
	t.Parallel()
	hotelID := uuid.New()
	tests := []struct {
		name       string
		hotelID    uuid.UUID
		roomNumber string
		roomType   domain.RoomType
		capacity   int
		wantErr    error
	}{
		{
			name:       "nil hotel id",
			hotelID:    uuid.Nil,
			roomNumber: "101",
			roomType:   domain.RoomTypeSingle,
			capacity:   1,
			wantErr:    domain.ErrInvalidRoom,
		},
		{
			name:       "empty room number",
			hotelID:    hotelID,
			roomNumber: "",
			roomType:   domain.RoomTypeSingle,
			capacity:   1,
			wantErr:    domain.ErrInvalidRoom,
		},
		{
			name:       "whitespace room number",
			hotelID:    hotelID,
			roomNumber: "   ",
			roomType:   domain.RoomTypeSingle,
			capacity:   1,
			wantErr:    domain.ErrInvalidRoom,
		},
		{
			name:       "invalid room type",
			hotelID:    hotelID,
			roomNumber: "101",
			roomType:   domain.RoomType("penthouse"),
			capacity:   2,
			wantErr:    domain.ErrInvalidRoom,
		},
		{
			name:       "zero capacity",
			hotelID:    hotelID,
			roomNumber: "101",
			roomType:   domain.RoomTypeSingle,
			capacity:   0,
			wantErr:    domain.ErrInvalidRoom,
		},
		{
			name:       "negative capacity",
			hotelID:    hotelID,
			roomNumber: "101",
			roomType:   domain.RoomTypeSingle,
			capacity:   -1,
			wantErr:    domain.ErrInvalidRoom,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := domain.NewRoom(tt.hotelID, tt.roomNumber, tt.roomType, tt.capacity)
			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGenerateDefaultRooms(t *testing.T) {
	t.Parallel()
	hotelID := uuid.New()

	rooms, err := domain.GenerateDefaultRooms(hotelID)
	require.NoError(t, err)
	require.Len(t, rooms, 6, "every hotel must have exactly 6 rooms")

	typeCounts := make(map[domain.RoomType]int)
	for _, r := range rooms {
		assert.Equal(t, hotelID, r.HotelID)
		assert.NotEmpty(t, r.RoomNumber)
		assert.Positive(t, r.Capacity)
		typeCounts[r.RoomType]++
	}

	assert.Equal(t, 2, typeCounts[domain.RoomTypeSingle], "must have 2 single rooms")
	assert.Equal(t, 2, typeCounts[domain.RoomTypeDouble], "must have 2 double rooms")
	assert.Equal(t, 2, typeCounts[domain.RoomTypeDeluxe], "must have 2 deluxe rooms")
}
