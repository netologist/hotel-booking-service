package httpapi

import (
	"github.com/netologist/hotel-booking/internal/adapter/httpapi/handler"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
)

type Server struct {
	*handler.HealthCheck
	*handler.ItemHandler
}

func NewServer(
	healthCheck *handler.HealthCheck,
	itemHandler *handler.ItemHandler,
) *Server {
	return &Server{
		HealthCheck: healthCheck,
		ItemHandler: itemHandler,
	}
}

var _ openapi.StrictServerInterface = (*Server)(nil)
