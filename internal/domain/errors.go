package domain

import "errors"

// ErrInvalidHotel indicates a domain invariant violation on a hotel entity.
var ErrInvalidHotel = errors.New("invalid hotel")

// ErrHotelAlreadyExists indicates a hotel with the same unique attribute already exists.
var ErrHotelAlreadyExists = errors.New("hotel already exists")

// ErrHotelNotFound indicates the requested hotel was not found.
var ErrHotelNotFound = errors.New("hotel not found")

// ErrInvalidRoom indicates a domain invariant violation on a room entity.
var ErrInvalidRoom = errors.New("invalid room")

// ErrRoomNotFound indicates the requested room was not found.
var ErrRoomNotFound = errors.New("room not found")
