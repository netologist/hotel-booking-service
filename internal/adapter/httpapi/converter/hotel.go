package converter

import (
	"github.com/netologist/hotel-booking/internal/domain"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
)

// goverter:converter
// goverter:output:file @cwd/internal/gen/goverter/httpapi/hotel.gen.go
// goverter:output:package github.com/netologist/hotel-booking/internal/gen/goverter/httpapi
// goverter:output:format function
// goverter:matchIgnoreCase
// goverter:skipCopySameType
type HotelConverter interface {
	// ToOpenAPIHotel converts a domain Hotel to an OpenAPI Hotel response.
	ToOpenAPIHotel(source domain.Hotel) openapi.Hotel

	// ToOpenAPIHotels converts a slice of domain Hotels to an OpenAPI Hotel response slice.
	ToOpenAPIHotels(source []domain.Hotel) []openapi.Hotel
}
