package httpapi

import (
	"github.com/netologist/hotel-booking/internal/adapter/httpapi/handler"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
)

type Server struct {
	*handler.HealthCheck
	*handler.HotelHandler
	*handler.AdminHandler
	*handler.RoomHandler
	*handler.BookingHandler
}

func NewServer(
	healthCheck *handler.HealthCheck,
	hotelHandler *handler.HotelHandler,
	adminHandler *handler.AdminHandler,
	roomHandler *handler.RoomHandler,
	bookingHandler *handler.BookingHandler,
) *Server {
	return &Server{
		HealthCheck:    healthCheck,
		HotelHandler:   hotelHandler,
		AdminHandler:   adminHandler,
		RoomHandler:    roomHandler,
		BookingHandler: bookingHandler,
	}
}

var _ openapi.StrictServerInterface = (*Server)(nil)
