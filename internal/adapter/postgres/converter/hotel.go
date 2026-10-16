package converter

import (
	"github.com/netologist/hotel-booking/internal/domain"
	"github.com/netologist/hotel-booking/internal/gen/sqlc"
)

// goverter:converter
// goverter:output:file @cwd/internal/gen/goverter/postgres/hotel.gen.go
// goverter:output:package github.com/netologist/hotel-booking/internal/gen/goverter/postgres
// goverter:output:format function
// goverter:matchIgnoreCase
// goverter:skipCopySameType
type HotelConverter interface {
	// ToDomainHotel converts a sqlc database Hotel to a domain Hotel entity.
	ToDomainHotel(source sqlc.Hotel) domain.Hotel

	// ToDomainHotels converts a slice of sqlc database Hotels to a domain Hotel entity slice.
	ToDomainHotels(source []sqlc.Hotel) []domain.Hotel

	// ToCreateHotelParams converts a domain Hotel entity to sqlc CreateHotelParams.
	ToCreateHotelParams(source domain.Hotel) sqlc.CreateHotelParams
}
