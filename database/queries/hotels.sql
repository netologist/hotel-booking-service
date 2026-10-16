-- name: CreateHotel :one
INSERT INTO hotels (id, name, address, city)
VALUES ($1, $2, $3, $4)
RETURNING id, name, address, city, created_at, updated_at;

-- name: SearchHotelsByName :many
SELECT id, name, address, city, created_at, updated_at FROM hotels
WHERE strpos(lower(name), lower(@name::text)) > 0
ORDER BY name;
