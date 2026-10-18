package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/netologist/hotel-booking/internal/domain"
)

type BookingService interface {
	BookRoom(
		ctx context.Context,
		roomID uuid.UUID,
		guestName, guestEmail string,
		checkInDate, checkOutDate string,
		numberOfGuests int,
	) (domain.Booking, error)
	GetBookingByReference(ctx context.Context, reference string) (domain.Booking, error)
}

type bookingService struct {
	roomRepo    domain.RoomRepository
	bookingRepo domain.BookingRepository
	notifier    BookingNotifier
}

var _ BookingService = (*bookingService)(nil)

func NewBookingService(
	roomRepo domain.RoomRepository,
	bookingRepo domain.BookingRepository,
	notifier BookingNotifier,
) *bookingService {
	return &bookingService{
		roomRepo:    roomRepo,
		bookingRepo: bookingRepo,
		notifier:    notifier,
	}
}

func (s *bookingService) BookRoom(
	ctx context.Context,
	roomID uuid.UUID,
	guestName, guestEmail string,
	checkInDate, checkOutDate string,
	numberOfGuests int,
) (domain.Booking, error) {
	checkIn, err := time.Parse("2006-01-02", checkInDate)
	if err != nil {
		return domain.Booking{}, fmt.Errorf("%w: invalid check-in date format", domain.ErrInvalidBooking)
	}

	checkOut, err := time.Parse("2006-01-02", checkOutDate)
	if err != nil {
		return domain.Booking{}, fmt.Errorf("%w: invalid check-out date format", domain.ErrInvalidBooking)
	}

	if !checkOut.After(checkIn) {
		return domain.Booking{}, fmt.Errorf(
			"%w: check-out date must be after check-in date",
			domain.ErrInvalidBooking,
		)
	}

	// Fetch room and verify existence
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return domain.Booking{}, fmt.Errorf("get room: %w", err)
	}

	// Validate domain invariants (capacity, name, email)
	booking, err := domain.NewBooking(room, guestName, guestEmail, checkIn, checkOut, numberOfGuests)
	if err != nil {
		return domain.Booking{}, err
	}

	// Check for conflicting overlapping bookings
	hasConflict, err := s.bookingRepo.CheckConflict(ctx, roomID, checkIn, checkOut)
	if err != nil {
		return domain.Booking{}, fmt.Errorf("check booking conflict: %w", err)
	}
	if hasConflict {
		return domain.Booking{}, domain.ErrRoomAlreadyBooked
	}

	// Persist booking
	createdBooking, err := s.bookingRepo.Create(ctx, booking)
	if err != nil {
		return domain.Booking{}, fmt.Errorf("create booking: %w", err)
	}

	// Asynchronously trigger confirmation email simulation
	if s.notifier != nil {
		s.notifier.SendConfirmationEmail(createdBooking)
	}

	return createdBooking, nil
}

func (s *bookingService) GetBookingByReference(ctx context.Context, reference string) (domain.Booking, error) {
	return s.bookingRepo.GetByReference(ctx, reference)
}
