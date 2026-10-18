package usecase_test

import (
	"context"
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

const (
	testHotelName    = "Grand Budapest Hotel"
	testHotelAddress = "100 Alpine Way"
	testHotelCity    = "Lutz"
)

func TestNewHotelService(t *testing.T) {
	mockRepo := domainmocks.NewMockHotelRepository(t)
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	service := usecase.NewHotelService(mockRepo, mockRoomRepo)
	require.NotNil(t, service)
}

func TestHotelService_CreateHotel_Success(t *testing.T) {
	mockRepo := domainmocks.NewMockHotelRepository(t)
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	service := usecase.NewHotelService(mockRepo, mockRoomRepo)
	ctx := t.Context()

	fixedTime := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)

	var capturedID uuid.UUID
	mockRepo.EXPECT().Create(ctx, mock.MatchedBy(func(h domain.Hotel) bool {
		capturedID = h.ID
		return h.Name == testHotelName &&
			h.Address == testHotelAddress &&
			h.City == testHotelCity &&
			h.ID != uuid.Nil
	})).RunAndReturn(func(_ context.Context, h domain.Hotel) (domain.Hotel, error) {
		return domain.Hotel{
			ID:        h.ID,
			Name:      h.Name,
			Address:   h.Address,
			City:      h.City,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		}, nil
	}).Once()

	mockRoomRepo.EXPECT().CreateBatch(ctx, mock.MatchedBy(func(rooms []domain.Room) bool {
		return len(rooms) == 6 && rooms[0].HotelID == capturedID
	})).Return(nil).Once()

	res, err := service.CreateHotel(ctx, testHotelName, testHotelAddress, testHotelCity)
	require.NoError(t, err)

	assert.Equal(t, capturedID, res.ID)
	assert.Equal(t, testHotelName, res.Name)
	assert.Equal(t, testHotelAddress, res.Address)
	assert.Equal(t, testHotelCity, res.City)
	assert.Len(t, res.Rooms, 6)
	assert.True(t, fixedTime.Equal(res.CreatedAt))
	assert.True(t, fixedTime.Equal(res.UpdatedAt))
}

func TestHotelService_CreateHotel_Error(t *testing.T) {
	mockRepo := domainmocks.NewMockHotelRepository(t)
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	service := usecase.NewHotelService(mockRepo, mockRoomRepo)
	ctx := t.Context()
	repoErr := errors.New("repository connection failed")

	mockRepo.EXPECT().Create(ctx, mock.MatchedBy(func(h domain.Hotel) bool {
		return h.Name == testHotelName &&
			h.Address == testHotelAddress &&
			h.City == testHotelCity &&
			h.ID != uuid.Nil
	})).Return(domain.Hotel{}, repoErr).Once()

	res, err := service.CreateHotel(ctx, testHotelName, testHotelAddress, testHotelCity)
	require.ErrorIs(t, err, repoErr)
	require.ErrorContains(t, err, "create hotel")
	assert.Equal(t, uuid.Nil, res.ID)
}

func TestHotelService_CreateHotel_DomainInvariantFailure(t *testing.T) {
	mockRepo := domainmocks.NewMockHotelRepository(t)
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	service := usecase.NewHotelService(mockRepo, mockRoomRepo)
	ctx := t.Context()

	res, err := service.CreateHotel(ctx, " ", testHotelAddress, testHotelCity)
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrInvalidHotel)
	assert.Equal(t, uuid.Nil, res.ID)
}

func TestHotelService_FindHotels_Success(t *testing.T) {
	mockRepo := domainmocks.NewMockHotelRepository(t)
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	service := usecase.NewHotelService(mockRepo, mockRoomRepo)
	ctx := t.Context()

	hotelID := uuid.New()
	expectedHotels := []domain.Hotel{
		{
			ID:      hotelID,
			Name:    testHotelName,
			Address: testHotelAddress,
			City:    testHotelCity,
		},
	}

	mockRepo.EXPECT().SearchByName(ctx, "Budapest").Return(expectedHotels, nil).Once()
	mockRoomRepo.EXPECT().GetByHotelID(ctx, hotelID).Return([]domain.Room{
		{ID: uuid.New(), HotelID: hotelID, RoomNumber: "101", RoomType: domain.RoomTypeSingle, Capacity: 1},
	}, nil).Once()

	res, err := service.FindHotels(ctx, "Budapest")
	require.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Len(t, res[0].Rooms, 1)
}

func TestHotelService_FindHotels_Error(t *testing.T) {
	mockRepo := domainmocks.NewMockHotelRepository(t)
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	service := usecase.NewHotelService(mockRepo, mockRoomRepo)
	ctx := t.Context()

	repoErr := errors.New("search failed")
	mockRepo.EXPECT().SearchByName(ctx, "Budapest").Return(nil, repoErr).Once()

	res, err := service.FindHotels(ctx, "Budapest")
	require.ErrorIs(t, err, repoErr)
	require.ErrorContains(t, err, "find hotels")
	assert.Nil(t, res)
}
