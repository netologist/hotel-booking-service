-- name: CreateBooking :one
INSERT INTO bookings (
    id, reference, room_id, hotel_id, guest_name, guest_email,
    check_in_date, check_out_date, number_of_guests
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, reference, room_id, hotel_id, guest_name, guest_email,
            check_in_date, check_out_date, number_of_guests, created_at;

-- name: CheckRoomBookingConflict :one
SELECT count(*) FROM bookings
WHERE room_id = @room_id
  AND check_in_date < @check_out_date::date
  AND check_out_date > @check_in_date::date;

-- name: GetBookingByReference :one
SELECT id, reference, room_id, hotel_id, guest_name, guest_email,
       check_in_date, check_out_date, number_of_guests, created_at
FROM bookings
WHERE reference = $1;

-- name: GetBookingDetailsByReference :one
SELECT 
    b.id,
    b.reference,
    b.guest_name,
    b.guest_email,
    b.check_in_date,
    b.check_out_date,
    b.number_of_guests,
    b.created_at,
    h.id AS hotel_id,
    h.name AS hotel_name,
    h.address AS hotel_address,
    h.city AS hotel_city,
    r.id AS room_id,
    r.room_number,
    r.room_type,
    r.capacity AS room_capacity
FROM bookings b
JOIN hotels h ON h.id = b.hotel_id
JOIN rooms r ON r.id = b.room_id
WHERE b.reference = $1;
