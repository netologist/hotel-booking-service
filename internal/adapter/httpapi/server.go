package httpapi

import (
	"github.com/netologist/hotel-booking/internal/adapter/httpapi/handler"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
)

type Server struct {
	*handler.HealthCheck
	*handler.HotelHandler
	*handler.AdminHandler
}

func NewServer(
	healthCheck *handler.HealthCheck,
	hotelHandler *handler.HotelHandler,
	adminHandler *handler.AdminHandler,
) *Server {
	return &Server{
		HealthCheck:  healthCheck,
		HotelHandler: hotelHandler,
		AdminHandler: adminHandler,
	}
}

var _ openapi.StrictServerInterface = (*Server)(nil)
