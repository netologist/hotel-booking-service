package handler_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/adapter/httpapi/handler"
	usecasemocks "github.com/netologist/hotel-booking/internal/gen/mocks/usecase"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
	"github.com/netologist/hotel-booking/internal/usecase"
)

func TestAdminHandler_ResetDatabase_Success(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockAdminService(t)
	h := handler.NewAdminHandler(mockService)
	ctx := t.Context()

	mockService.EXPECT().Reset(ctx).Return(nil).Once()

	resp, err := h.ResetDatabase(ctx, openapi.ResetDatabaseRequestObject{})
	require.NoError(t, err)

	jsonResp, ok := resp.(openapi.ResetDatabase200JSONResponse)
	require.True(t, ok)
	assert.Equal(t, "Database reset successfully", jsonResp.Message)
}

func TestAdminHandler_ResetDatabase_Error(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockAdminService(t)
	h := handler.NewAdminHandler(mockService)
	ctx := t.Context()

	resetErr := errors.New("reset failed")
	mockService.EXPECT().Reset(ctx).Return(resetErr).Once()

	_, err := h.ResetDatabase(ctx, openapi.ResetDatabaseRequestObject{})
	require.ErrorIs(t, err, resetErr)
}

func TestAdminHandler_SeedDatabase_Success(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockAdminService(t)
	h := handler.NewAdminHandler(mockService)
	ctx := t.Context()

	mockService.EXPECT().Seed(ctx).Return(usecase.SeedResult{
		HotelsCount: 2,
		RoomsCount:  12,
	}, nil).Once()

	resp, err := h.SeedDatabase(ctx, openapi.SeedDatabaseRequestObject{})
	require.NoError(t, err)

	jsonResp, ok := resp.(openapi.SeedDatabase200JSONResponse)
	require.True(t, ok)
	assert.Equal(t, "Database seeded successfully", jsonResp.Message)
	assert.Equal(t, 2, jsonResp.HotelsCount)
	assert.Equal(t, 12, jsonResp.RoomsCount)
}

func TestAdminHandler_SeedDatabase_Error(t *testing.T) {
	t.Parallel()
	mockService := usecasemocks.NewMockAdminService(t)
	h := handler.NewAdminHandler(mockService)
	ctx := t.Context()

	seedErr := errors.New("seed failed")
	mockService.EXPECT().Seed(ctx).Return(usecase.SeedResult{}, seedErr).Once()

	_, err := h.SeedDatabase(ctx, openapi.SeedDatabaseRequestObject{})
	require.ErrorIs(t, err, seedErr)
}
