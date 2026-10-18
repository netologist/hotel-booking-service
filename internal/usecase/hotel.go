package usecase

import (
	"context"
	"fmt"

	"github.com/netologist/hotel-booking/internal/domain"
)

// HotelService is the inbound port for hotel business operations.
type HotelService interface {
	CreateHotel(ctx context.Context, name, address, city string) (domain.Hotel, error)
	FindHotels(ctx context.Context, name string) ([]domain.Hotel, error)
}

// hotelService implements HotelService.
type hotelService struct {
	hotelRepo domain.HotelRepository
	roomRepo  domain.RoomRepository
}

var _ HotelService = (*hotelService)(nil)

// NewHotelService constructor returns a concrete *hotelService.
func NewHotelService(hotelRepo domain.HotelRepository, roomRepo domain.RoomRepository) *hotelService {
	return &hotelService{
		hotelRepo: hotelRepo,
		roomRepo:  roomRepo,
	}
}

func (u *hotelService) CreateHotel(ctx context.Context, name, address, city string) (domain.Hotel, error) {
	hotel, err := domain.NewHotel(name, address, city)
	if err != nil {
		return domain.Hotel{}, fmt.Errorf("create hotel: %w", err)
	}

	entity, err := u.hotelRepo.Create(ctx, hotel)
	if err != nil {
		return domain.Hotel{}, fmt.Errorf("create hotel: %w", err)
	}

	rooms, err := domain.GenerateDefaultRooms(entity.ID)
	if err != nil {
		return domain.Hotel{}, fmt.Errorf("generate default rooms: %w", err)
	}

	if err := u.roomRepo.CreateBatch(ctx, rooms); err != nil {
		return domain.Hotel{}, fmt.Errorf("create default rooms: %w", err)
	}

	entity.Rooms = rooms
	return entity, nil
}

func (u *hotelService) FindHotels(ctx context.Context, name string) ([]domain.Hotel, error) {
	hotels, err := u.hotelRepo.SearchByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("find hotels: %w", err)
	}

	for i := range hotels {
		rooms, err := u.roomRepo.GetByHotelID(ctx, hotels[i].ID)
		if err == nil && len(rooms) > 0 {
			hotels[i].Rooms = rooms
		}
	}

	return hotels, nil
}
