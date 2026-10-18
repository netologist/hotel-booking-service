package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/netologist/hotel-booking/internal/domain"
)

type RoomService interface {
	FindAvailableRooms(
		ctx context.Context,
		hotelID *uuid.UUID,
		checkInDate, checkOutDate string,
		guests int,
	) ([]domain.AvailableRoom, error)
}

type roomService struct {
	roomRepo domain.RoomRepository
}

var _ RoomService = (*roomService)(nil)

func NewRoomService(roomRepo domain.RoomRepository) *roomService {
	return &roomService{
		roomRepo: roomRepo,
	}
}

func (s *roomService) FindAvailableRooms(
	ctx context.Context,
	hotelID *uuid.UUID,
	checkInDate, checkOutDate string,
	guests int,
) ([]domain.AvailableRoom, error) {
	if guests < 1 {
		return nil, fmt.Errorf("%w: number of guests must be at least 1", domain.ErrInvalidRoom)
	}

	checkIn, err := time.Parse("2006-01-02", checkInDate)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid check-in date format, expected YYYY-MM-DD", domain.ErrInvalidRoom)
	}

	checkOut, err := time.Parse("2006-01-02", checkOutDate)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid check-out date format, expected YYYY-MM-DD", domain.ErrInvalidRoom)
	}

	if !checkOut.After(checkIn) {
		return nil, fmt.Errorf("%w: check-out date must be after check-in date", domain.ErrInvalidRoom)
	}

	rooms, err := s.roomRepo.FindAvailableRooms(ctx, hotelID, checkIn, checkOut, guests)
	if err != nil {
		return nil, fmt.Errorf("find available rooms: %w", err)
	}

	return rooms, nil
}
