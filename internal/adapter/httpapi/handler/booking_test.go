package handler_test

import (
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

const testGuestJohnDoe = "John Doe"

func TestBookingHandler_BookRoom_Success(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockBookingService(t)
	h := handler.NewBookingHandler(mockService)
	ctx := t.Context()

	roomID := uuid.New()
	bookingID := uuid.New()
	hotelID := uuid.New()
	checkIn := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	checkOut := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)

	expectedBooking := domain.Booking{
		ID:             bookingID,
		Reference:      "BK-TEST-1234",
		RoomID:         roomID,
		HotelID:        hotelID,
		GuestName:      testGuestJohnDoe,
		GuestEmail:     "john@example.com",
		CheckInDate:    checkIn,
		CheckOutDate:   checkOut,
		NumberOfGuests: 2,
		CreatedAt:      time.Now(),
	}

	mockService.EXPECT().BookRoom(ctx, roomID, testGuestJohnDoe, "john@example.com", "2026-11-01", "2026-11-03", 2).
		Return(expectedBooking, nil).Once()

	req := openapi.BookRoomRequestObject{
		Body: &openapi.BookRoomRequest{
			RoomId:         roomID,
			CheckInDate:    openapi_types.Date{Time: checkIn},
			CheckOutDate:   openapi_types.Date{Time: checkOut},
			NumberOfGuests: 2,
			GuestName:      testGuestJohnDoe,
			GuestEmail:     openapi_types.Email("john@example.com"),
		},
	}

	resp, err := h.BookRoom(ctx, req)
	require.NoError(t, err)

	jsonResp, ok := resp.(openapi.BookRoom201JSONResponse)
	require.True(t, ok)
	assert.Equal(t, "BK-TEST-1234", jsonResp.Reference)
	assert.Equal(t, testGuestJohnDoe, jsonResp.GuestName)
}

func TestBookingHandler_BookRoom_Conflict(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockBookingService(t)
	h := handler.NewBookingHandler(mockService)
	ctx := t.Context()

	roomID := uuid.New()
	mockService.EXPECT().
		BookRoom(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(domain.Booking{}, domain.ErrRoomAlreadyBooked).Once()

	req := openapi.BookRoomRequestObject{
		Body: &openapi.BookRoomRequest{
			RoomId:         roomID,
			CheckInDate:    openapi_types.Date{Time: time.Now()},
			CheckOutDate:   openapi_types.Date{Time: time.Now().Add(24 * time.Hour)},
			NumberOfGuests: 1,
			GuestName:      testGuestJohnDoe,
			GuestEmail:     openapi_types.Email("john@example.com"),
		},
	}

	resp, err := h.BookRoom(ctx, req)
	require.NoError(t, err)

	_, ok := resp.(openapi.BookRoom409ApplicationProblemPlusJSONResponse)
	require.True(t, ok)
}

func TestBookingHandler_BookRoom_CapacityExceeded(t *testing.T) { //nolint:dupl // test payload variations
	t.Parallel()
	mockService := usecasemocks.NewMockBookingService(t)
	h := handler.NewBookingHandler(mockService)
	ctx := t.Context()

	mockService.EXPECT().
		BookRoom(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(domain.Booking{}, domain.ErrCapacityExceeded).Once()

	req := openapi.BookRoomRequestObject{
		Body: &openapi.BookRoomRequest{
			RoomId:         uuid.New(),
			CheckInDate:    openapi_types.Date{Time: time.Now()},
			CheckOutDate:   openapi_types.Date{Time: time.Now().Add(24 * time.Hour)},
			NumberOfGuests: 5,
			GuestName:      testGuestJohnDoe,
			GuestEmail:     openapi_types.Email("john@example.com"),
		},
	}

	resp, err := h.BookRoom(ctx, req)
	require.NoError(t, err)

	_, ok := resp.(openapi.BookRoom400ApplicationProblemPlusJSONResponse)
	require.True(t, ok)
}

func TestBookingHandler_BookRoom_RoomNotFound(t *testing.T) { //nolint:dupl // test payload variations
	t.Parallel()
	mockService := usecasemocks.NewMockBookingService(t)
	h := handler.NewBookingHandler(mockService)
	ctx := t.Context()

	mockService.EXPECT().
		BookRoom(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(domain.Booking{}, domain.ErrRoomNotFound).Once()

	req := openapi.BookRoomRequestObject{
		Body: &openapi.BookRoomRequest{
			RoomId:         uuid.New(),
			CheckInDate:    openapi_types.Date{Time: time.Now()},
			CheckOutDate:   openapi_types.Date{Time: time.Now().Add(24 * time.Hour)},
			NumberOfGuests: 1,
			GuestName:      testGuestJohnDoe,
			GuestEmail:     openapi_types.Email("john@example.com"),
		},
	}

	resp, err := h.BookRoom(ctx, req)
	require.NoError(t, err)

	_, ok := resp.(openapi.BookRoom404ApplicationProblemPlusJSONResponse)
	require.True(t, ok)
}

func TestBookingHandler_GetBookingByReference_Success(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockBookingService(t)
	h := handler.NewBookingHandler(mockService)
	ctx := t.Context()

	expected := domain.BookingDetails{
		Reference:    "BK-REF-100",
		GuestName:    testGuestJohnDoe,
		GuestEmail:   "john@example.com",
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(24 * time.Hour),
		HotelName:    "Grand Budapest Hotel",
		RoomNumber:   "201",
		RoomType:     domain.RoomTypeDouble,
		RoomCapacity: 2,
	}

	mockService.EXPECT().GetBookingByReference(ctx, "BK-REF-100").Return(expected, nil).Once()

	resp, err := h.GetBookingByReference(ctx, openapi.GetBookingByReferenceRequestObject{
		BookingReference: "BK-REF-100",
	})
	require.NoError(t, err)

	jsonResp, ok := resp.(openapi.GetBookingByReference200JSONResponse)
	require.True(t, ok)
	assert.Equal(t, "BK-REF-100", jsonResp.Reference)
	assert.Equal(t, "Grand Budapest Hotel", jsonResp.Hotel.Name)
}

func TestBookingHandler_GetBookingByReference_NotFound(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockBookingService(t)
	h := handler.NewBookingHandler(mockService)
	ctx := t.Context()

	mockService.EXPECT().GetBookingByReference(ctx, "NONEXISTENT").
		Return(domain.BookingDetails{}, domain.ErrBookingNotFound).Once()

	resp, err := h.GetBookingByReference(ctx, openapi.GetBookingByReferenceRequestObject{
		BookingReference: "NONEXISTENT",
	})
	require.NoError(t, err)

	_, ok := resp.(openapi.GetBookingByReference404ApplicationProblemPlusJSONResponse)
	require.True(t, ok)
}

func TestBookingHandler_GetBookingByReference_Blank(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockBookingService(t)
	h := handler.NewBookingHandler(mockService)
	ctx := t.Context()

	resp, err := h.GetBookingByReference(ctx, openapi.GetBookingByReferenceRequestObject{
		BookingReference: "   ",
	})
	require.NoError(t, err)

	_, ok := resp.(openapi.GetBookingByReference400ApplicationProblemPlusJSONResponse)
	require.True(t, ok)
}
