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

-- name: FindAvailableRooms :many
SELECT 
    r.id, 
    r.hotel_id, 
    h.name AS hotel_name, 
    r.room_number, 
    r.room_type, 
    r.capacity, 
    r.created_at, 
    r.updated_at
FROM rooms r
JOIN hotels h ON h.id = r.hotel_id
WHERE r.capacity >= @guests::int
  AND (sqlc.narg('hotel_id')::uuid IS NULL OR r.hotel_id = sqlc.narg('hotel_id')::uuid)
  AND r.id NOT IN (
      SELECT b.room_id
      FROM bookings b
      WHERE b.check_in_date < @check_out_date::date
        AND b.check_out_date > @check_in_date::date
  )
ORDER BY h.name, r.room_number;
