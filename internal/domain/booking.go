package domain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidBooking    = errors.New("invalid booking")
	ErrCapacityExceeded  = errors.New("number of guests exceeds room capacity")
	ErrRoomAlreadyBooked = errors.New("room already booked for selected dates")
	ErrBookingNotFound   = errors.New("booking not found")
)

const (
	MinGuestNameLength    = 2
	bookingRefRandomBytes = 4
)

type Booking struct {
	ID             uuid.UUID
	Reference      string
	RoomID         uuid.UUID
	HotelID        uuid.UUID
	GuestName      string
	GuestEmail     string
	CheckInDate    time.Time
	CheckOutDate   time.Time
	NumberOfGuests int
	CreatedAt      time.Time
}

func GenerateBookingReference() string {
	b := make([]byte, bookingRefRandomBytes)
	_, _ = rand.Read(b)
	return fmt.Sprintf("BK-%s-%s", time.Now().Format("20060102"), strings.ToUpper(hex.EncodeToString(b)))
}

func NewBooking(
	room Room,
	guestName, guestEmail string,
	checkIn, checkOut time.Time,
	numberOfGuests int,
) (Booking, error) {
	if room.ID == uuid.Nil {
		return Booking{}, fmt.Errorf("%w: room is required", ErrInvalidBooking)
	}
	guestName = strings.TrimSpace(guestName)
	if len(guestName) < MinGuestNameLength {
		return Booking{}, fmt.Errorf("%w: guest name must be at least 2 characters", ErrInvalidBooking)
	}
	guestEmail = strings.TrimSpace(guestEmail)
	if !strings.Contains(guestEmail, "@") {
		return Booking{}, fmt.Errorf("%w: invalid guest email address", ErrInvalidBooking)
	}
	if !checkOut.After(checkIn) {
		return Booking{}, fmt.Errorf("%w: check-out date must be after check-in date", ErrInvalidBooking)
	}
	if numberOfGuests < 1 {
		return Booking{}, fmt.Errorf("%w: number of guests must be at least 1", ErrInvalidBooking)
	}
	if numberOfGuests > room.Capacity {
		return Booking{}, fmt.Errorf(
			"%w: room capacity is %d, requested %d",
			ErrCapacityExceeded,
			room.Capacity,
			numberOfGuests,
		)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Booking{}, fmt.Errorf("generate booking id: %w", err)
	}

	return Booking{
		ID:             id,
		Reference:      GenerateBookingReference(),
		RoomID:         room.ID,
		HotelID:        room.HotelID,
		GuestName:      guestName,
		GuestEmail:     guestEmail,
		CheckInDate:    checkIn,
		CheckOutDate:   checkOut,
		NumberOfGuests: numberOfGuests,
		CreatedAt:      time.Now(),
	}, nil
}

type BookingRepository interface {
	Create(ctx context.Context, booking Booking) (Booking, error)
	CheckConflict(ctx context.Context, roomID uuid.UUID, checkIn, checkOut time.Time) (bool, error)
	GetByReference(ctx context.Context, reference string) (Booking, error)
}
