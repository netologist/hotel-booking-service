package domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RoomType string

const (
	RoomTypeSingle RoomType = "single"
	RoomTypeDouble RoomType = "double"
	RoomTypeDeluxe RoomType = "deluxe"

	SingleCapacity = 1
	DoubleCapacity = 2
	DeluxeCapacity = 3
)

func (t RoomType) IsValid() bool {
	switch t {
	case RoomTypeSingle, RoomTypeDouble, RoomTypeDeluxe:
		return true
	default:
		return false
	}
}

func (t RoomType) DefaultCapacity() int {
	switch t {
	case RoomTypeSingle:
		return SingleCapacity
	case RoomTypeDouble:
		return DoubleCapacity
	case RoomTypeDeluxe:
		return DeluxeCapacity
	default:
		return 0
	}
}

type Room struct {
	ID         uuid.UUID
	HotelID    uuid.UUID
	RoomNumber string
	RoomType   RoomType
	Capacity   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewRoom(hotelID uuid.UUID, roomNumber string, roomType RoomType, capacity int) (Room, error) {
	if hotelID == uuid.Nil {
		return Room{}, fmt.Errorf("%w: hotel id is required", ErrInvalidRoom)
	}
	roomNumber = strings.TrimSpace(roomNumber)
	if roomNumber == "" {
		return Room{}, fmt.Errorf("%w: room number is required", ErrInvalidRoom)
	}
	if !roomType.IsValid() {
		return Room{}, fmt.Errorf("%w: invalid room type %q", ErrInvalidRoom, roomType)
	}
	if capacity <= 0 {
		return Room{}, fmt.Errorf("%w: capacity must be greater than 0", ErrInvalidRoom)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return Room{}, fmt.Errorf("generate id: %w", err)
	}

	return Room{
		ID:         id,
		HotelID:    hotelID,
		RoomNumber: roomNumber,
		RoomType:   roomType,
		Capacity:   capacity,
	}, nil
}

// GenerateDefaultRooms creates the standard 6 rooms for a hotel: 2 single, 2 double, 2 deluxe.
func GenerateDefaultRooms(hotelID uuid.UUID) ([]Room, error) {
	specs := []struct {
		number   string
		roomType RoomType
		capacity int
	}{
		{"101", RoomTypeSingle, 1},
		{"102", RoomTypeSingle, 1},
		{"201", RoomTypeDouble, 2},
		{"202", RoomTypeDouble, 2},
		{"301", RoomTypeDeluxe, 3},
		{"302", RoomTypeDeluxe, 3},
	}

	rooms := make([]Room, 0, len(specs))
	for _, s := range specs {
		r, err := NewRoom(hotelID, s.number, s.roomType, s.capacity)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}
	return rooms, nil
}

type RoomRepository interface {
	CreateBatch(ctx context.Context, rooms []Room) error
	GetByHotelID(ctx context.Context, hotelID uuid.UUID) ([]Room, error)
	GetByID(ctx context.Context, id uuid.UUID) (Room, error)
}
