package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/netologist/hotel-booking/internal/adapter/httpapi/converter"
	"github.com/netologist/hotel-booking/internal/adapter/httpapi/problem"
	"github.com/netologist/hotel-booking/internal/domain"
	goverter_httpapi "github.com/netologist/hotel-booking/internal/gen/goverter/httpapi"
	"github.com/netologist/hotel-booking/internal/gen/openapi"
	"github.com/netologist/hotel-booking/internal/usecase"
)

type BookingHandler struct {
	bookingService usecase.BookingService
}

func NewBookingHandler(bookingService usecase.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) BookRoom(
	ctx context.Context,
	request openapi.BookRoomRequestObject,
) (openapi.BookRoomResponseObject, error) {
	req := request.Body
	checkIn := req.CheckInDate.String()
	checkOut := req.CheckOutDate.String()

	booking, err := h.bookingService.BookRoom(
		ctx,
		req.RoomId,
		req.GuestName,
		string(req.GuestEmail),
		checkIn,
		checkOut,
		req.NumberOfGuests,
	)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidBooking) || errors.Is(err, domain.ErrCapacityExceeded) {
			prob := problem.BadRequest(err.Error(), "/api/v1/bookings", "INVALID_BOOKING_REQUEST", nil)
			return openapi.BookRoom400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: openapi.BadRequestApplicationProblemPlusJSONResponse(prob),
			}, nil
		}
		if errors.Is(err, domain.ErrRoomNotFound) {
			prob := problem.NotFound(err.Error(), "/api/v1/bookings", "ROOM_NOT_FOUND")
			return openapi.BookRoom404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: openapi.NotFoundApplicationProblemPlusJSONResponse(prob),
			}, nil
		}
		if errors.Is(err, domain.ErrRoomAlreadyBooked) {
			prob := problem.Conflict(err.Error(), "/api/v1/bookings", "ROOM_ALREADY_BOOKED")
			return openapi.BookRoom409ApplicationProblemPlusJSONResponse{
				ConflictApplicationProblemPlusJSONResponse: openapi.ConflictApplicationProblemPlusJSONResponse(prob),
			}, nil
		}
		return nil, err
	}

	return openapi.BookRoom201JSONResponse(goverter_httpapi.ToOpenAPIBooking(booking)), nil
}

func (h *BookingHandler) GetBookingByReference(
	ctx context.Context,
	request openapi.GetBookingByReferenceRequestObject,
) (openapi.GetBookingByReferenceResponseObject, error) {
	ref := strings.TrimSpace(request.BookingReference)
	if ref == "" {
		prob := problem.BadRequest("booking reference is required", "/api/v1/bookings", "INVALID_BOOKING_REFERENCE", nil)
		return openapi.GetBookingByReference400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: openapi.BadRequestApplicationProblemPlusJSONResponse(prob),
		}, nil
	}

	details, err := h.bookingService.GetBookingByReference(ctx, ref)
	if err != nil {
		if errors.Is(err, domain.ErrBookingNotFound) {
			prob := problem.NotFound("booking not found", "/api/v1/bookings/"+ref, "BOOKING_NOT_FOUND")
			return openapi.GetBookingByReference404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: openapi.NotFoundApplicationProblemPlusJSONResponse(prob),
			}, nil
		}
		return nil, err
	}

	return openapi.GetBookingByReference200JSONResponse(converter.ToOpenAPIBookingDetails(details)), nil
}
