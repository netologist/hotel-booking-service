package postgres

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/domain"
	"github.com/netologist/hotel-booking/internal/gen/sqlc"
)

const testRoomTypeSingle = "single"

func TestNewRoomRepository(t *testing.T) {
	repo := NewRoomRepository(nil)
	require.NotNil(t, repo)
}

func TestRoomRepository_CreateBatch_Success(t *testing.T) {
	q := &mockQuerier{}
	repo := NewRoomRepositoryWithQuerier(q)
	ctx := t.Context()

	hotelID := uuid.New()
	roomID := uuid.New()
	rooms := []domain.Room{
		{
			ID:         roomID,
			HotelID:    hotelID,
			RoomNumber: "101",
			RoomType:   domain.RoomTypeSingle,
			Capacity:   1,
		},
	}

	q.On("CreateRoom", ctx, mock.MatchedBy(func(p sqlc.CreateRoomParams) bool {
		return p.ID == roomID && p.HotelID == hotelID && p.RoomNumber == "101" && p.Capacity == 1
	})).Return(sqlc.Room{
		ID:         roomID,
		HotelID:    hotelID,
		RoomNumber: "101",
		RoomType:   testRoomTypeSingle,
		Capacity:   1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil).Once()

	err := repo.CreateBatch(ctx, rooms)
	require.NoError(t, err)
	q.AssertExpectations(t)
}

func TestRoomRepository_CreateBatch_Error(t *testing.T) {
	q := &mockQuerier{}
	repo := NewRoomRepositoryWithQuerier(q)
	ctx := t.Context()

	rooms := []domain.Room{
		{
			ID:         uuid.New(),
			HotelID:    uuid.New(),
			RoomNumber: "101",
			RoomType:   domain.RoomTypeSingle,
			Capacity:   1,
		},
	}

	dbErr := errors.New("db error")
	q.On("CreateRoom", ctx, mock.Anything).Return(sqlc.Room{}, dbErr).Once()

	err := repo.CreateBatch(ctx, rooms)
	require.ErrorIs(t, err, dbErr)
	q.AssertExpectations(t)
}

func TestRoomRepository_GetByHotelID(t *testing.T) {
	q := &mockQuerier{}
	repo := NewRoomRepositoryWithQuerier(q)
	ctx := t.Context()

	hotelID := uuid.New()
	q.On("GetRoomsByHotelID", ctx, hotelID).Return([]sqlc.Room{
		{
			ID:         uuid.New(),
			HotelID:    hotelID,
			RoomNumber: "101",
			RoomType:   testRoomTypeSingle,
			Capacity:   1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}, nil).Once()

	rooms, err := repo.GetByHotelID(ctx, hotelID)
	require.NoError(t, err)
	require.Len(t, rooms, 1)
	assert.Equal(t, "101", rooms[0].RoomNumber)
	assert.Equal(t, domain.RoomTypeSingle, rooms[0].RoomType)
	assert.Equal(t, 1, rooms[0].Capacity)
	q.AssertExpectations(t)
}

func TestRoomRepository_GetByID_NotFound(t *testing.T) {
	q := &mockQuerier{}
	repo := NewRoomRepositoryWithQuerier(q)
	ctx := t.Context()

	roomID := uuid.New()
	q.On("GetRoomByID", ctx, roomID).Return(sqlc.Room{}, pgx.ErrNoRows).Once()

	_, err := repo.GetByID(ctx, roomID)
	require.ErrorIs(t, err, domain.ErrRoomNotFound)
	q.AssertExpectations(t)
}

func TestRoomRepository_FindAvailableRooms(t *testing.T) {
	q := &mockQuerier{}
	repo := NewRoomRepositoryWithQuerier(q)
	ctx := t.Context()

	hotelID := uuid.New()
	checkIn := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	checkOut := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)

	q.On("FindAvailableRooms", ctx, mock.Anything).Return([]sqlc.FindAvailableRoomsRow{
		{
			ID:         uuid.New(),
			HotelID:    hotelID,
			HotelName:  "Grand Budapest Hotel",
			RoomNumber: "101",
			RoomType:   testRoomTypeSingle,
			Capacity:   1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}, nil).Once()

	rooms, err := repo.FindAvailableRooms(ctx, &hotelID, checkIn, checkOut, 1)
	require.NoError(t, err)
	require.Len(t, rooms, 1)
	assert.Equal(t, "101", rooms[0].RoomNumber)
	assert.Equal(t, "Grand Budapest Hotel", rooms[0].HotelName)
	assert.Equal(t, domain.RoomTypeSingle, rooms[0].RoomType)
	q.AssertExpectations(t)
}
