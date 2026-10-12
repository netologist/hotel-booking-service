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

func SetupTestRouterWithMock(t *testing.T) (http.Handler, *usecasemocks.MockItemService) {
	mockItemService := usecasemocks.NewMockItemService(t)
	server := httpapi.NewServer(
		handler.NewHealthCheck(),
		handler.NewItemHandler(mockItemService),
	)
	return httpapi.NewRouter(server), mockItemService
}

func SetupTestRouterWithHealth(t *testing.T, health *handler.HealthCheck) http.Handler {
	mockItemService := usecasemocks.NewMockItemService(t)
	server := httpapi.NewServer(
		health,
		handler.NewItemHandler(mockItemService),
	)
	return httpapi.NewRouter(server)
}
