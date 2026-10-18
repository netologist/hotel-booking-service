-- +goose Up
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY,
    reference VARCHAR(64) NOT NULL UNIQUE,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    hotel_id UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    guest_name VARCHAR(255) NOT NULL,
    guest_email VARCHAR(255) NOT NULL,
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    number_of_guests INT NOT NULL CHECK (number_of_guests > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_booking_dates CHECK (check_out_date > check_in_date)
);

CREATE INDEX IF NOT EXISTS idx_bookings_reference ON bookings (reference);
CREATE INDEX IF NOT EXISTS idx_bookings_room_dates ON bookings (room_id, check_in_date, check_out_date);
CREATE INDEX IF NOT EXISTS idx_bookings_hotel_id ON bookings (hotel_id);

-- +goose Down
DROP TABLE IF EXISTS bookings;
