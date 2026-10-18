package handler

import (
	"context"
	"errors"

	"github.com/netologist/hotel-booking/internal/adapter/httpapi/problem"
	"github.com/netologist/hotel-booking/internal/domain"
	goverter_httpapi "github.com/netologist/hotel-booking/internal/gen/goverter/httpapi"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
	"github.com/netologist/hotel-booking/internal/usecase"
)

type RoomHandler struct {
	roomService usecase.RoomService
}

func NewRoomHandler(roomService usecase.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

func (h *RoomHandler) FindAvailableRooms(
	ctx context.Context,
	request openapi.FindAvailableRoomsRequestObject,
) (openapi.FindAvailableRoomsResponseObject, error) {
	checkIn := request.Params.CheckInDate.String()
	checkOut := request.Params.CheckOutDate.String()
	guests := request.Params.Guests

	rooms, err := h.roomService.FindAvailableRooms(ctx, request.Params.HotelId, checkIn, checkOut, guests)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidRoom) {
			prob := problem.BadRequest(err.Error(), "/api/v1/rooms/available", "INVALID_QUERY_PARAMETER", nil)
			return openapi.FindAvailableRooms400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: openapi.BadRequestApplicationProblemPlusJSONResponse(prob),
			}, nil
		}
		return nil, err
	}

	responseRooms := goverter_httpapi.ToOpenAPIAvailableRooms(rooms)
	if responseRooms == nil {
		responseRooms = []openapi.AvailableRoom{}
	}

	return openapi.FindAvailableRooms200JSONResponse{
		Rooms: responseRooms,
	}, nil
}
