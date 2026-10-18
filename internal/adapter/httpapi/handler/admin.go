package handler

import (
	"context"

	"github.com/netologist/hotel-booking/internal/gen/openapi"
	"github.com/netologist/hotel-booking/internal/usecase"
)

type AdminHandler struct {
	adminService usecase.AdminService
}

func NewAdminHandler(adminService usecase.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) ResetDatabase(
	ctx context.Context,
	_ openapi.ResetDatabaseRequestObject,
) (openapi.ResetDatabaseResponseObject, error) {
	if err := h.adminService.Reset(ctx); err != nil {
		return nil, err
	}
	return openapi.ResetDatabase200JSONResponse{
		Message: "Database reset successfully",
	}, nil
}

func (h *AdminHandler) SeedDatabase(
	ctx context.Context,
	_ openapi.SeedDatabaseRequestObject,
) (openapi.SeedDatabaseResponseObject, error) {
	result, err := h.adminService.Seed(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.SeedDatabase200JSONResponse{
		Message:     "Database seeded successfully",
		HotelsCount: result.HotelsCount,
		RoomsCount:  result.RoomsCount,
	}, nil
}
