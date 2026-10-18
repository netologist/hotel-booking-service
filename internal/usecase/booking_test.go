package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/domain"
	domainmocks "github.com/netologist/hotel-booking/internal/gen/mocks/domain"
	usecasemocks "github.com/netologist/hotel-booking/internal/gen/mocks/usecase"
	"github.com/netologist/hotel-booking/internal/usecase"
)

const (
	testBookingCheckIn  = "2026-11-01"
	testBookingCheckOut = "2026-11-03"
	testGuestName       = "John Doe"
	testGuestEmail      = "john@example.com"
)

func TestBookingService_BookRoom_Success(t *testing.T) {
	t.Parallel()
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	mockBookingRepo := domainmocks.NewMockBookingRepository(t)
	mockNotifier := usecasemocks.NewMockBookingNotifier(t)

	service := usecase.NewBookingService(mockRoomRepo, mockBookingRepo, mockNotifier)
	ctx := t.Context()

	roomID := uuid.New()
	hotelID := uuid.New()
	room := domain.Room{
		ID:         roomID,
		HotelID:    hotelID,
		RoomNumber: "201",
		RoomType:   domain.RoomTypeDouble,
		Capacity:   2,
	}

	mockRoomRepo.EXPECT().GetByID(ctx, roomID).Return(room, nil).Once()
	mockBookingRepo.EXPECT().CheckConflict(ctx, roomID, mock.Anything, mock.Anything).Return(false, nil).Once()
	mockBookingRepo.EXPECT().Create(ctx, mock.MatchedBy(func(b domain.Booking) bool {
		return b.RoomID == roomID && b.GuestName == testGuestName && b.NumberOfGuests == 2
	})).RunAndReturn(func(_ context.Context, b domain.Booking) (domain.Booking, error) {
		return b, nil
	}).Once()

	mockNotifier.EXPECT().SendConfirmationEmail(mock.MatchedBy(func(b domain.Booking) bool {
		return b.RoomID == roomID
	})).Once()

	booking, err := service.BookRoom(
		ctx, roomID, testGuestName, testGuestEmail, testBookingCheckIn, testBookingCheckOut, 2,
	)
	require.NoError(t, err)
	assert.Equal(t, roomID, booking.RoomID)
	assert.Equal(t, testGuestName, booking.GuestName)
}

func TestBookingService_BookRoom_Conflict(t *testing.T) {
	t.Parallel()
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	mockBookingRepo := domainmocks.NewMockBookingRepository(t)
	mockNotifier := usecasemocks.NewMockBookingNotifier(t)

	service := usecase.NewBookingService(mockRoomRepo, mockBookingRepo, mockNotifier)
	ctx := t.Context()

	roomID := uuid.New()
	room := domain.Room{
		ID:         roomID,
		HotelID:    uuid.New(),
		RoomNumber: "201",
		RoomType:   domain.RoomTypeDouble,
		Capacity:   2,
	}

	mockRoomRepo.EXPECT().GetByID(ctx, roomID).Return(room, nil).Once()
	mockBookingRepo.EXPECT().CheckConflict(ctx, roomID, mock.Anything, mock.Anything).Return(true, nil).Once()

	_, err := service.BookRoom(ctx, roomID, testGuestName, testGuestEmail, testBookingCheckIn, testBookingCheckOut, 2)
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrRoomAlreadyBooked)
}

func TestBookingService_BookRoom_CapacityExceeded(t *testing.T) {
	t.Parallel()
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	mockBookingRepo := domainmocks.NewMockBookingRepository(t)

	service := usecase.NewBookingService(mockRoomRepo, mockBookingRepo, nil)
	ctx := t.Context()

	roomID := uuid.New()
	room := domain.Room{
		ID:         roomID,
		HotelID:    uuid.New(),
		RoomNumber: "101",
		RoomType:   domain.RoomTypeSingle,
		Capacity:   1,
	}

	mockRoomRepo.EXPECT().GetByID(ctx, roomID).Return(room, nil).Once()

	_, err := service.BookRoom(ctx, roomID, testGuestName, testGuestEmail, testBookingCheckIn, testBookingCheckOut, 2)
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrCapacityExceeded)
}

func TestBookingService_BookRoom_InvalidDates(t *testing.T) {
	t.Parallel()
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	mockBookingRepo := domainmocks.NewMockBookingRepository(t)

	service := usecase.NewBookingService(mockRoomRepo, mockBookingRepo, nil)
	ctx := t.Context()

	roomID := uuid.New()

	// check-out before check-in
	_, err := service.BookRoom(ctx, roomID, testGuestName, testGuestEmail, "2026-11-05", "2026-11-01", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrInvalidBooking)
}

func TestBookingService_BookRoom_RoomNotFound(t *testing.T) {
	t.Parallel()
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	mockBookingRepo := domainmocks.NewMockBookingRepository(t)

	service := usecase.NewBookingService(mockRoomRepo, mockBookingRepo, nil)
	ctx := t.Context()

	roomID := uuid.New()
	mockRoomRepo.EXPECT().GetByID(ctx, roomID).Return(domain.Room{}, domain.ErrRoomNotFound).Once()

	_, err := service.BookRoom(ctx, roomID, testGuestName, testGuestEmail, testBookingCheckIn, testBookingCheckOut, 1)
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrRoomNotFound)
}

func TestBookingService_GetBookingByReference_Success(t *testing.T) {
	t.Parallel()
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	mockBookingRepo := domainmocks.NewMockBookingRepository(t)

	service := usecase.NewBookingService(mockRoomRepo, mockBookingRepo, nil)
	ctx := t.Context()

	expected := domain.BookingDetails{
		Reference: "BK-REF-1",
		GuestName: "John Doe",
	}
	mockBookingRepo.EXPECT().GetDetailsByReference(ctx, "BK-REF-1").Return(expected, nil).Once()

	details, err := service.GetBookingByReference(ctx, "BK-REF-1")
	require.NoError(t, err)
	assert.Equal(t, expected, details)
}

func TestBookingService_GetBookingByReference_Error(t *testing.T) {
	t.Parallel()
	mockRoomRepo := domainmocks.NewMockRoomRepository(t)
	mockBookingRepo := domainmocks.NewMockBookingRepository(t)

	service := usecase.NewBookingService(mockRoomRepo, mockBookingRepo, nil)
	ctx := t.Context()

	mockBookingRepo.EXPECT().GetDetailsByReference(ctx, "BK-REF-1").
		Return(domain.BookingDetails{}, domain.ErrBookingNotFound).Once()

	_, err := service.GetBookingByReference(ctx, "BK-REF-1")
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrBookingNotFound)
}

func TestAsyncEmailNotifier(t *testing.T) {
	t.Parallel()
	notifier := usecase.NewAsyncEmailNotifier(10 * time.Millisecond)
	booking := domain.Booking{
		Reference:    "BK-TEST",
		GuestName:    "John",
		GuestEmail:   "john@example.com",
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(24 * time.Hour),
	}

	// Should not panic or block
	notifier.SendConfirmationEmail(booking)
	time.Sleep(20 * time.Millisecond)
}

func TestAsyncEmailNotifier_Drain(t *testing.T) {
	t.Parallel()
	notifier := usecase.NewAsyncEmailNotifier(10 * time.Millisecond)
	booking := domain.Booking{
		Reference:    "BK-DRAIN-TEST",
		GuestName:    "John",
		GuestEmail:   "john@example.com",
		CheckInDate:  time.Now(),
		CheckOutDate: time.Now().Add(24 * time.Hour),
	}

	notifier.SendConfirmationEmail(booking)

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	err := notifier.Drain(ctx)
	require.NoError(t, err)
}
