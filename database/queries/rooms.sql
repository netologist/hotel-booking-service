-- name: CreateRoom :one
INSERT INTO rooms (id, hotel_id, room_number, room_type, capacity)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, hotel_id, room_number, room_type, capacity, created_at, updated_at;

-- name: GetRoomsByHotelID :many
SELECT id, hotel_id, room_number, room_type, capacity, created_at, updated_at
FROM rooms
WHERE hotel_id = $1
ORDER BY room_number;

-- name: GetRoomByID :one
SELECT id, hotel_id, room_number, room_type, capacity, created_at, updated_at
FROM rooms
WHERE id = $1;
