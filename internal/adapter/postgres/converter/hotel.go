package converter

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/netologist/hotel-booking/internal/domain"
	"github.com/netologist/hotel-booking/internal/gen/sqlc"
)

// goverter:converter
// goverter:output:file @cwd/internal/gen/goverter/postgres/hotel.gen.go
// goverter:output:package github.com/netologist/hotel-booking/internal/gen/goverter/postgres
// goverter:output:format function
// goverter:matchIgnoreCase
// goverter:skipCopySameType
// goverter:extend IntToInt32 Int32ToInt DateToTime TimeToDate
type HotelConverter interface {
	// ToDomainHotel converts a sqlc database Hotel to a domain Hotel entity.
	// goverter:ignore Rooms
	ToDomainHotel(source sqlc.Hotel) domain.Hotel

	// ToDomainHotels converts a slice of sqlc database Hotels to a domain Hotel entity slice.
	ToDomainHotels(source []sqlc.Hotel) []domain.Hotel

	// ToCreateHotelParams converts a domain Hotel entity to sqlc CreateHotelParams.
	ToCreateHotelParams(source domain.Hotel) sqlc.CreateHotelParams

	// ToDomainRoom converts a sqlc database Room to a domain Room entity.
	ToDomainRoom(source sqlc.Room) domain.Room

	// ToDomainRooms converts a slice of sqlc database Rooms to a domain Room entity slice.
	ToDomainRooms(source []sqlc.Room) []domain.Room

	// ToCreateRoomParams converts a domain Room entity to sqlc CreateRoomParams.
	ToCreateRoomParams(source domain.Room) sqlc.CreateRoomParams

	// ToDomainAvailableRoom converts a sqlc FindAvailableRoomsRow to a domain AvailableRoom entity.
	ToDomainAvailableRoom(source sqlc.FindAvailableRoomsRow) domain.AvailableRoom

	// ToDomainAvailableRooms converts a slice of sqlc FindAvailableRoomsRow to a domain AvailableRoom slice.
	ToDomainAvailableRooms(source []sqlc.FindAvailableRoomsRow) []domain.AvailableRoom

	// ToDomainBooking converts a sqlc Booking to a domain Booking entity.
	ToDomainBooking(source sqlc.Booking) domain.Booking

	// ToCreateBookingParams converts a domain Booking entity to sqlc CreateBookingParams.
	ToCreateBookingParams(source domain.Booking) sqlc.CreateBookingParams
}

func IntToInt32(source int) int32 {
	return int32(source) //nolint:gosec // room capacity is a small bounded integer
}

func Int32ToInt(source int32) int {
	return int(source)
}

func DateToTime(source pgtype.Date) time.Time {
	return source.Time
}

func TimeToDate(source time.Time) pgtype.Date {
	return pgtype.Date{Time: source, Valid: true}
}
