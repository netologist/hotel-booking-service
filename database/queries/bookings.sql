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
