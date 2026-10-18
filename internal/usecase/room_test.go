package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/domain"
	domainmocks "github.com/netologist/hotel-booking/internal/gen/mocks/domain"
	"github.com/netologist/hotel-booking/internal/usecase"
)

const testDate20261101 = "2026-11-01"

func TestRoomService_FindAvailableRooms_Success(t *testing.T) {
	t.Parallel()
	mockRepo := domainmocks.NewMockRoomRepository(t)
	service := usecase.NewRoomService(mockRepo)
	ctx := t.Context()

	hotelID := uuid.New()
	expectedRooms := []domain.AvailableRoom{
		{
			ID:         uuid.New(),
			HotelID:    hotelID,
			HotelName:  "Grand Budapest Hotel",
			RoomNumber: "101",
			RoomType:   domain.RoomTypeSingle,
			Capacity:   1,
		},
	}

	mockRepo.EXPECT().FindAvailableRooms(
		ctx,
		&hotelID,
		time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC),
		1,
	).Return(expectedRooms, nil).Once()

	rooms, err := service.FindAvailableRooms(ctx, &hotelID, testDate20261101, "2026-11-03", 1)
	require.NoError(t, err)
	assert.Equal(t, expectedRooms, rooms)
}

func TestRoomService_FindAvailableRooms_ValidationErrors(t *testing.T) {
	t.Parallel()
	mockRepo := domainmocks.NewMockRoomRepository(t)
	service := usecase.NewRoomService(mockRepo)
	ctx := t.Context()

	hotelID := uuid.New()

	tests := []struct {
		name        string
		guests      int
		checkIn     string
		checkOut    string
		errContains string
	}{
		{
			name:        "guests less than 1",
			guests:      0,
			checkIn:     testDate20261101,
			checkOut:    "2026-11-03",
			errContains: "number of guests must be at least 1",
		},
		{
			name:        "invalid check-in format",
			guests:      1,
			checkIn:     "invalid-date",
			checkOut:    "2026-11-03",
			errContains: "invalid check-in date format",
		},
		{
			name:        "invalid check-out format",
			guests:      1,
			checkIn:     testDate20261101,
			checkOut:    "invalid-date",
			errContains: "invalid check-out date format",
		},
		{
			name:        "check-out before check-in",
			guests:      1,
			checkIn:     "2026-11-05",
			checkOut:    testDate20261101,
			errContains: "check-out date must be after check-in date",
		},
		{
			name:        "check-out same as check-in",
			guests:      1,
			checkIn:     testDate20261101,
			checkOut:    testDate20261101,
			errContains: "check-out date must be after check-in date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := service.FindAvailableRooms(ctx, &hotelID, tt.checkIn, tt.checkOut, tt.guests)
			require.Error(t, err)
			require.ErrorIs(t, err, domain.ErrInvalidRoom)
			assert.Contains(t, err.Error(), tt.errContains)
		})
	}
}

func TestRoomService_FindAvailableRooms_RepoError(t *testing.T) {
	t.Parallel()
	mockRepo := domainmocks.NewMockRoomRepository(t)
	service := usecase.NewRoomService(mockRepo)
	ctx := t.Context()

	dbErr := errors.New("database failure")
	mockRepo.EXPECT().FindAvailableRooms(ctx, (*uuid.UUID)(nil), mock.Anything, mock.Anything, 2).Return(nil, dbErr).Once()

	_, err := service.FindAvailableRooms(ctx, nil, testDate20261101, "2026-11-03", 2)
	require.ErrorIs(t, err, dbErr)
}
