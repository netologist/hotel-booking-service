package usecase

import (
	"context"
	"fmt"

	"github.com/netologist/hotel-booking/internal/domain"
)

type SeedResult struct {
	HotelsCount int
	RoomsCount  int
}

type AdminService interface {
	Reset(ctx context.Context) error
	Seed(ctx context.Context) (SeedResult, error)
}

type adminService struct {
	hotelRepo domain.HotelRepository
	roomRepo  domain.RoomRepository
	adminRepo domain.AdminRepository
}

var _ AdminService = (*adminService)(nil)

func NewAdminService(
	hotelRepo domain.HotelRepository,
	roomRepo domain.RoomRepository,
	adminRepo domain.AdminRepository,
) *adminService {
	return &adminService{
		hotelRepo: hotelRepo,
		roomRepo:  roomRepo,
		adminRepo: adminRepo,
	}
}

func (a *adminService) Reset(ctx context.Context) error {
	if err := a.adminRepo.TruncateAll(ctx); err != nil {
		return fmt.Errorf("reset database: %w", err)
	}
	return nil
}

func (a *adminService) Seed(ctx context.Context) (SeedResult, error) {
	if err := a.adminRepo.TruncateAll(ctx); err != nil {
		return SeedResult{}, fmt.Errorf("seed reset: %w", err)
	}

	seedHotels := []struct {
		name    string
		address string
		city    string
	}{
		{"Grand Budapest Hotel", "100 Alpine Way", "Lutz"},
		{"The Continental", "1 Wall Street", "New York"},
	}

	totalRooms := 0
	for _, h := range seedHotels {
		hotel, err := domain.NewHotel(h.name, h.address, h.city)
		if err != nil {
			return SeedResult{}, fmt.Errorf("seed new hotel: %w", err)
		}

		createdHotel, err := a.hotelRepo.Create(ctx, hotel)
		if err != nil {
			return SeedResult{}, fmt.Errorf("seed create hotel: %w", err)
		}

		rooms, err := domain.GenerateDefaultRooms(createdHotel.ID)
		if err != nil {
			return SeedResult{}, fmt.Errorf("seed generate rooms: %w", err)
		}

		if err := a.roomRepo.CreateBatch(ctx, rooms); err != nil {
			return SeedResult{}, fmt.Errorf("seed create rooms: %w", err)
		}

		totalRooms += len(rooms)
	}

	return SeedResult{
		HotelsCount: len(seedHotels),
		RoomsCount:  totalRooms,
	}, nil
}
