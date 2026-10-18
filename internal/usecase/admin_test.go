package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/domain"
	domainmocks "github.com/netologist/hotel-booking/internal/gen/mocks/domain"
	"github.com/netologist/hotel-booking/internal/usecase"
)

func TestAdminService_Reset(t *testing.T) {
	t.Parallel()
	mockAdminRepo := domainmocks.NewMockAdminRepository(t)
	mockHotelRepo := domainmocks.NewMockHotelRepository(t)
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)

	service := usecase.NewAdminService(mockHotelRepo, mockRoomRepo, mockAdminRepo)
	ctx := t.Context()

	mockAdminRepo.EXPECT().TruncateAll(ctx).Return(nil).Once()

	err := service.Reset(ctx)
	require.NoError(t, err)
}

func TestAdminService_Reset_Error(t *testing.T) {
	t.Parallel()
	mockAdminRepo := domainmocks.NewMockAdminRepository(t)
	mockHotelRepo := domainmocks.NewMockHotelRepository(t)
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)

	service := usecase.NewAdminService(mockHotelRepo, mockRoomRepo, mockAdminRepo)
	ctx := t.Context()

	truncateErr := errors.New("db truncate error")
	mockAdminRepo.EXPECT().TruncateAll(ctx).Return(truncateErr).Once()

	err := service.Reset(ctx)
	require.ErrorIs(t, err, truncateErr)
}

func TestAdminService_Seed(t *testing.T) {
	t.Parallel()
	mockAdminRepo := domainmocks.NewMockAdminRepository(t)
	mockHotelRepo := domainmocks.NewMockHotelRepository(t)
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)

	service := usecase.NewAdminService(mockHotelRepo, mockRoomRepo, mockAdminRepo)
	ctx := t.Context()

	mockAdminRepo.EXPECT().TruncateAll(ctx).Return(nil).Once()

	mockHotelRepo.EXPECT().Create(ctx, mock.Anything).RunAndReturn(
		func(_ context.Context, h domain.Hotel) (domain.Hotel, error) {
			return domain.Hotel{
				ID:      uuid.New(),
				Name:    h.Name,
				Address: h.Address,
				City:    h.City,
			}, nil
		}).Times(2)

	mockRoomRepo.EXPECT().CreateBatch(ctx, mock.MatchedBy(func(rooms []domain.Room) bool {
		return len(rooms) == 6
	})).Return(nil).Times(2)

	result, err := service.Seed(ctx)
	require.NoError(t, err)
	assert.Equal(t, 2, result.HotelsCount)
	assert.Equal(t, 12, result.RoomsCount)
}
