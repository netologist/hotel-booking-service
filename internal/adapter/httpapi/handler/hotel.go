package handler

import (
	"context"
	"errors"

	"github.com/netologist/hotel-booking/internal/adapter/httpapi/problem"
	"github.com/netologist/hotel-booking/internal/domain"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
	"github.com/netologist/hotel-booking/internal/usecase"
)

func NewHotelHandler(hotelService usecase.HotelService) *HotelHandler {
	return &HotelHandler{
		hotelService: hotelService,
	}
}

type HotelHandler struct {
	hotelService usecase.HotelService
}

func (h *HotelHandler) CreateHotel(
	ctx context.Context,
	request openapi.CreateHotelRequestObject,
) (openapi.CreateHotelResponseObject, error) {
	model, err := h.hotelService.CreateHotel(ctx, request.Body.Name, request.Body.Address, request.Body.City)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidHotel) {
			prob := problem.BadRequest(err.Error(), "/api/v1/hotels", "INVALID_HOTEL", nil)
			return openapi.CreateHotel400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: openapi.BadRequestApplicationProblemPlusJSONResponse(prob),
			}, nil
		}
		if errors.Is(err, domain.ErrHotelAlreadyExists) {
			prob := problem.Conflict(domain.ErrHotelAlreadyExists.Error(), "/api/v1/hotels", "HOTEL_ALREADY_EXISTS")
			return openapi.CreateHotel409ApplicationProblemPlusJSONResponse{
				ConflictApplicationProblemPlusJSONResponse: openapi.ConflictApplicationProblemPlusJSONResponse(prob),
			}, nil
		}
		return nil, err
	}

	return openapi.CreateHotel201JSONResponse(toHotelResponse(model)), nil
}

// mapper methods

func toHotelResponse(output domain.Hotel) openapi.Hotel {
	return openapi.Hotel{
		Id:        output.ID,
		Name:      output.Name,
		City:      output.City,
		Address:   output.Address,
		CreatedAt: output.CreatedAt,
		UpdatedAt: output.UpdatedAt,
	}
}
