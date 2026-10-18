package handler_test

import (
	"net/http"
	"testing"

	"github.com/netologist/hotel-booking/internal/adapter/httpapi"
	"github.com/netologist/hotel-booking/internal/adapter/httpapi/handler"
	usecasemocks "github.com/netologist/hotel-booking/internal/gen/mocks/usecase"
)

func SetupTestRouter(t *testing.T) http.Handler {
	router, _ := SetupTestRouterWithMock(t)
	return router
}

func SetupTestRouterWithMock(t *testing.T) (http.Handler, *usecasemocks.MockHotelService) {
	mockHotelService := usecasemocks.NewMockHotelService(t)
	mockAdminService := usecasemocks.NewMockAdminService(t)
	mockRoomService := usecasemocks.NewMockRoomService(t)
	mockBookingService := usecasemocks.NewMockBookingService(t)
	server := httpapi.NewServer(
		handler.NewHealthCheck(),
		handler.NewHotelHandler(mockHotelService),
		handler.NewAdminHandler(mockAdminService),
		handler.NewRoomHandler(mockRoomService),
		handler.NewBookingHandler(mockBookingService),
	)
	return httpapi.NewRouter(server), mockHotelService
}

func SetupTestRouterWithHealth(t *testing.T, health *handler.HealthCheck) http.Handler {
	mockHotelService := usecasemocks.NewMockHotelService(t)
	mockAdminService := usecasemocks.NewMockAdminService(t)
	mockRoomService := usecasemocks.NewMockRoomService(t)
	mockBookingService := usecasemocks.NewMockBookingService(t)
	server := httpapi.NewServer(
		health,
		handler.NewHotelHandler(mockHotelService),
		handler.NewAdminHandler(mockAdminService),
		handler.NewRoomHandler(mockRoomService),
		handler.NewBookingHandler(mockBookingService),
	)
	return httpapi.NewRouter(server)
}
