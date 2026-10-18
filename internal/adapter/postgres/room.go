package postgres

import (
	"context"
	"errors"

	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/netologist/hotel-booking/internal/domain"
	goverter_postgres "github.com/netologist/hotel-booking/internal/gen/goverter/postgres"
	"github.com/netologist/hotel-booking/internal/gen/sqlc"
)

type roomRepository struct {
	queries sqlc.Querier
}

var _ domain.RoomRepository = (*roomRepository)(nil)

func NewRoomRepository(pool *pgxpool.Pool) *roomRepository {
	return NewRoomRepositoryWithQuerier(sqlc.New(pool))
}

func NewRoomRepositoryWithQuerier(queries sqlc.Querier) *roomRepository {
	return &roomRepository{queries: queries}
}

func (r *roomRepository) CreateBatch(ctx context.Context, rooms []domain.Room) error {
	for _, room := range rooms {
		params := goverter_postgres.ToCreateRoomParams(room)
		if _, err := r.queries.CreateRoom(ctx, params); err != nil {
			return err
		}
	}
	return nil
}

func (r *roomRepository) GetByHotelID(ctx context.Context, hotelID uuid.UUID) ([]domain.Room, error) {
	rows, err := r.queries.GetRoomsByHotelID(ctx, hotelID)
	if err != nil {
		return nil, err
	}
	return goverter_postgres.ToDomainRooms(rows), nil
}

func (r *roomRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Room, error) {
	row, err := r.queries.GetRoomByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Room{}, domain.ErrRoomNotFound
		}
		return domain.Room{}, err
	}
	return goverter_postgres.ToDomainRoom(row), nil
}

func (r *roomRepository) FindAvailableRooms(
	ctx context.Context,
	hotelID *uuid.UUID,
	checkIn, checkOut time.Time,
	guests int,
) ([]domain.AvailableRoom, error) {
	var pgHotelID pgtype.UUID
	if hotelID != nil && *hotelID != uuid.Nil {
		pgHotelID = pgtype.UUID{Bytes: *hotelID, Valid: true}
	}

	params := sqlc.FindAvailableRoomsParams{
		Guests:       int32(guests), //nolint:gosec // guests count is a small bounded integer
		HotelID:      pgHotelID,
		CheckInDate:  pgtype.Date{Time: checkIn, Valid: true},
		CheckOutDate: pgtype.Date{Time: checkOut, Valid: true},
	}

	rows, err := r.queries.FindAvailableRooms(ctx, params)
	if err != nil {
		return nil, err
	}

	return goverter_postgres.ToDomainAvailableRooms(rows), nil
}
