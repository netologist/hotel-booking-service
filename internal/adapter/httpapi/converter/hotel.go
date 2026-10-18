package converter

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/netologist/hotel-booking/internal/domain"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
)

// goverter:converter
// goverter:output:file @cwd/internal/gen/goverter/httpapi/hotel.gen.go
// goverter:output:package github.com/netologist/hotel-booking/internal/gen/goverter/httpapi
// goverter:output:format function
// goverter:matchIgnoreCase
// goverter:skipCopySameType
// goverter:extend ToOpenAPIRoomType TimeToOpenAPIDate
type HotelConverter interface {
	// ToOpenAPIHotel converts a domain Hotel to an OpenAPI Hotel response.
	ToOpenAPIHotel(source domain.Hotel) openapi.Hotel

	// ToOpenAPIHotels converts a slice of domain Hotels to an OpenAPI Hotel response slice.
	ToOpenAPIHotels(source []domain.Hotel) []openapi.Hotel

	// ToOpenAPIRoom converts a domain Room to an OpenAPI Room response.
	ToOpenAPIRoom(source domain.Room) openapi.Room

	// ToOpenAPIRooms converts a slice of domain Rooms to an OpenAPI Room response slice.
	ToOpenAPIRooms(source []domain.Room) []openapi.Room

	// ToOpenAPIAvailableRoom converts a domain AvailableRoom to an OpenAPI AvailableRoom response.
	ToOpenAPIAvailableRoom(source domain.AvailableRoom) openapi.AvailableRoom

	// ToOpenAPIAvailableRooms converts a slice of domain AvailableRooms to an OpenAPI AvailableRoom response slice.
	ToOpenAPIAvailableRooms(source []domain.AvailableRoom) []openapi.AvailableRoom

	// ToOpenAPIBooking converts a domain Booking to an OpenAPI Booking response.
	ToOpenAPIBooking(source domain.Booking) openapi.Booking
}

func ToOpenAPIRoomType(source domain.RoomType) openapi.RoomType {
	return openapi.RoomType(source)
}

func TimeToOpenAPIDate(source time.Time) openapi_types.Date {
	return openapi_types.Date{Time: source}
}
