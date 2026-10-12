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
	service := usecase.NewHotelService(mockRepo)
	require.NotNil(t, service)
}

func TestHotelService_CreateHotel_Success(t *testing.T) {
	mockRepo := domainmocks.NewMockHotelRepository(t)
	service := usecase.NewHotelService(mockRepo)
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

	res, err := service.CreateHotel(ctx, testHotelName, testHotelAddress, testHotelCity)
	require.NoError(t, err)

	assert.Equal(t, capturedID, res.ID)
	assert.Equal(t, testHotelName, res.Name)
	assert.Equal(t, testHotelAddress, res.Address)
	assert.Equal(t, testHotelCity, res.City)
	assert.True(t, fixedTime.Equal(res.CreatedAt))
	assert.True(t, fixedTime.Equal(res.UpdatedAt))
}

func TestHotelService_CreateHotel_Error(t *testing.T) {
	mockRepo := domainmocks.NewMockHotelRepository(t)
	service := usecase.NewHotelService(mockRepo)
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
	service := usecase.NewHotelService(mockRepo)
	ctx := t.Context()

	res, err := service.CreateHotel(ctx, " ", testHotelAddress, testHotelCity)
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrInvalidHotel)
	assert.Equal(t, uuid.Nil, res.ID)
}
