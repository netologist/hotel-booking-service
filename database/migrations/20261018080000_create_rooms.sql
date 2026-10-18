-- +goose Up
CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY,
    hotel_id UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    room_number VARCHAR(20) NOT NULL,
    room_type VARCHAR(20) NOT NULL CHECK (room_type IN ('single', 'double', 'deluxe')),
    capacity INT NOT NULL CHECK (capacity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_hotel_room_number UNIQUE (hotel_id, room_number)
);

CREATE INDEX IF NOT EXISTS idx_rooms_hotel_id ON rooms (hotel_id);
CREATE INDEX IF NOT EXISTS idx_rooms_capacity ON rooms (capacity);

-- +goose Down
DROP TABLE IF EXISTS rooms;
