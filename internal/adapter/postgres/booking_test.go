package postgres

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/domain"
	"github.com/netologist/hotel-booking/internal/gen/sqlc"
)

const (
	testBookingRef = "BK-123456"
	testGuestName  = "John Doe"
	testGuestEmail = "john@example.com"
)

func TestNewBookingRepository(t *testing.T) {
	repo := NewBookingRepository(nil)
	require.NotNil(t, repo)
}

func TestBookingRepository_Create_Success(t *testing.T) {
	q := &mockQuerier{}
	repo := NewBookingRepositoryWithQuerier(nil, q)
	ctx := t.Context()

	bookingID := uuid.New()
	roomID := uuid.New()
	hotelID := uuid.New()
	checkIn := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	checkOut := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)

	b := domain.Booking{
		ID:             bookingID,
		Reference:      testBookingRef,
		RoomID:         roomID,
		HotelID:        hotelID,
		GuestName:      testGuestName,
		GuestEmail:     testGuestEmail,
		CheckInDate:    checkIn,
		CheckOutDate:   checkOut,
		NumberOfGuests: 2,
	}

	q.On("CreateBooking", ctx, mock.MatchedBy(func(p sqlc.CreateBookingParams) bool {
		return p.ID == bookingID && p.Reference == testBookingRef
	})).Return(sqlc.Booking{
		ID:             bookingID,
		Reference:      testBookingRef,
		RoomID:         roomID,
		HotelID:        hotelID,
		GuestName:      testGuestName,
		GuestEmail:     testGuestEmail,
		CheckInDate:    pgtype.Date{Time: checkIn, Valid: true},
		CheckOutDate:   pgtype.Date{Time: checkOut, Valid: true},
		NumberOfGuests: 2,
		CreatedAt:      time.Now(),
	}, nil).Once()

	created, err := repo.Create(ctx, b)
	require.NoError(t, err)
	assert.Equal(t, "BK-123456", created.Reference)
	q.AssertExpectations(t)
}

func TestBookingRepository_CheckConflict(t *testing.T) {
	q := &mockQuerier{}
	repo := NewBookingRepositoryWithQuerier(nil, q)
	ctx := t.Context()

	roomID := uuid.New()
	checkIn := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	checkOut := time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC)

	q.On("CheckRoomBookingConflict", ctx, mock.Anything).Return(int64(1), nil).Once()

	hasConflict, err := repo.CheckConflict(ctx, roomID, checkIn, checkOut)
	require.NoError(t, err)
	assert.True(t, hasConflict)
	q.AssertExpectations(t)
}

func TestBookingRepository_GetByReference_NotFound(t *testing.T) {
	q := &mockQuerier{}
	repo := NewBookingRepositoryWithQuerier(nil, q)
	ctx := t.Context()

	q.On("GetBookingByReference", ctx, "NONEXISTENT").Return(sqlc.Booking{}, pgx.ErrNoRows).Once()

	_, err := repo.GetByReference(ctx, "NONEXISTENT")
	require.ErrorIs(t, err, domain.ErrBookingNotFound)
	q.AssertExpectations(t)
}

func TestBookingRepository_GetDetailsByReference_Success(t *testing.T) {
	q := &mockQuerier{}
	repo := NewBookingRepositoryWithQuerier(nil, q)
	ctx := t.Context()

	q.On("GetBookingDetailsByReference", ctx, testBookingRef).Return(sqlc.GetBookingDetailsByReferenceRow{
		ID:             uuid.New(),
		Reference:      testBookingRef,
		GuestName:      testGuestName,
		GuestEmail:     testGuestEmail,
		CheckInDate:    pgtype.Date{Time: time.Now(), Valid: true},
		CheckOutDate:   pgtype.Date{Time: time.Now().Add(24 * time.Hour), Valid: true},
		NumberOfGuests: 2,
		CreatedAt:      time.Now(),
		HotelID:        uuid.New(),
		HotelName:      "Grand Budapest Hotel",
		HotelAddress:   "100 Alpine Way",
		HotelCity:      "Lutz",
		RoomID:         uuid.New(),
		RoomNumber:     "201",
		RoomType:       "double",
		RoomCapacity:   2,
	}, nil).Once()

	details, err := repo.GetDetailsByReference(ctx, testBookingRef)
	require.NoError(t, err)
	assert.Equal(t, testBookingRef, details.Reference)
	assert.Equal(t, "Grand Budapest Hotel", details.HotelName)
	q.AssertExpectations(t)
}

func TestBookingRepository_GetDetailsByReference_NotFound(t *testing.T) {
	q := &mockQuerier{}
	repo := NewBookingRepositoryWithQuerier(nil, q)
	ctx := t.Context()

	q.On("GetBookingDetailsByReference", ctx, "NONEXISTENT").
		Return(sqlc.GetBookingDetailsByReferenceRow{}, pgx.ErrNoRows).Once()

	_, err := repo.GetDetailsByReference(ctx, "NONEXISTENT")
	require.ErrorIs(t, err, domain.ErrBookingNotFound)
	q.AssertExpectations(t)
}
