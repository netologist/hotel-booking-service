package handler_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/adapter/httpapi/handler"
	"github.com/netologist/hotel-booking/internal/domain"
	usecasemocks "github.com/netologist/hotel-booking/internal/gen/mocks/usecase"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
)

func TestRoomHandler_FindAvailableRooms_Success(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockRoomService(t)
	h := handler.NewRoomHandler(mockService)
	ctx := t.Context()

	hotelID := uuid.New()
	mockService.EXPECT().FindAvailableRooms(ctx, &hotelID, "2026-11-01", "2026-11-03", 2).Return([]domain.AvailableRoom{
		{
			ID:         uuid.New(),
			HotelID:    hotelID,
			HotelName:  "Grand Budapest Hotel",
			RoomNumber: "201",
			RoomType:   domain.RoomTypeDouble,
			Capacity:   2,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}, nil).Once()

	req := openapi.FindAvailableRoomsRequestObject{
		Params: openapi.FindAvailableRoomsParams{
			CheckInDate:  openapi_types.Date{Time: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)},
			CheckOutDate: openapi_types.Date{Time: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)},
			Guests:       2,
			HotelId:      &hotelID,
		},
	}

	resp, err := h.FindAvailableRooms(ctx, req)
	require.NoError(t, err)

	jsonResp, ok := resp.(openapi.FindAvailableRooms200JSONResponse)
	require.True(t, ok)
	require.Len(t, jsonResp.Rooms, 1)
	assert.Equal(t, "201", jsonResp.Rooms[0].RoomNumber)
}

func TestRoomHandler_FindAvailableRooms_ValidationError(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockRoomService(t)
	h := handler.NewRoomHandler(mockService)
	ctx := t.Context()

	mockService.EXPECT().FindAvailableRooms(ctx, mock.Anything, mock.Anything, mock.Anything, 0).
		Return(nil, domain.ErrInvalidRoom).Once()

	req := openapi.FindAvailableRoomsRequestObject{
		Params: openapi.FindAvailableRoomsParams{
			CheckInDate:  openapi_types.Date{Time: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)},
			CheckOutDate: openapi_types.Date{Time: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)},
			Guests:       0,
		},
	}

	resp, err := h.FindAvailableRooms(ctx, req)
	require.NoError(t, err)

	_, ok := resp.(openapi.FindAvailableRooms400ApplicationProblemPlusJSONResponse)
	require.True(t, ok)
}

func TestRoomHandler_FindAvailableRooms_ServiceError(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockRoomService(t)
	h := handler.NewRoomHandler(mockService)
	ctx := t.Context()

	serviceErr := errors.New("database failure")
	mockService.EXPECT().FindAvailableRooms(ctx, mock.Anything, mock.Anything, mock.Anything, 1).
		Return(nil, serviceErr).Once()

	req := openapi.FindAvailableRoomsRequestObject{
		Params: openapi.FindAvailableRoomsParams{
			CheckInDate:  openapi_types.Date{Time: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)},
			CheckOutDate: openapi_types.Date{Time: time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)},
			Guests:       1,
		},
	}

	_, err := h.FindAvailableRooms(ctx, req)
	require.ErrorIs(t, err, serviceErr)
}
