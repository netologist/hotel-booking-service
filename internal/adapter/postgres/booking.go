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

type bookingRepository struct {
	pool    *pgxpool.Pool
	queries sqlc.Querier
}

var _ domain.BookingRepository = (*bookingRepository)(nil)

func NewBookingRepository(pool *pgxpool.Pool) *bookingRepository {
	return NewBookingRepositoryWithQuerier(pool, sqlc.New(pool))
}

func NewBookingRepositoryWithQuerier(pool *pgxpool.Pool, queries sqlc.Querier) *bookingRepository {
	return &bookingRepository{
		pool:    pool,
		queries: queries,
	}
}

func (b *bookingRepository) Create(ctx context.Context, booking domain.Booking) (domain.Booking, error) {
	params := goverter_postgres.ToCreateBookingParams(booking)
	row, err := b.queries.CreateBooking(ctx, params)
	if err != nil {
		return domain.Booking{}, err
	}
	return goverter_postgres.ToDomainBooking(row), nil
}

func (b *bookingRepository) CheckConflict(
	ctx context.Context,
	roomID uuid.UUID,
	checkIn, checkOut time.Time,
) (bool, error) {
	count, err := b.queries.CheckRoomBookingConflict(ctx, sqlc.CheckRoomBookingConflictParams{
		RoomID:       roomID,
		CheckInDate:  pgtype.Date{Time: checkIn, Valid: true},
		CheckOutDate: pgtype.Date{Time: checkOut, Valid: true},
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (b *bookingRepository) GetByReference(ctx context.Context, reference string) (domain.Booking, error) {
	row, err := b.queries.GetBookingByReference(ctx, reference)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Booking{}, domain.ErrBookingNotFound
		}
		return domain.Booking{}, err
	}
	return goverter_postgres.ToDomainBooking(row), nil
}

func (b *bookingRepository) GetDetailsByReference(
	ctx context.Context,
	reference string,
) (domain.BookingDetails, error) {
	row, err := b.queries.GetBookingDetailsByReference(ctx, reference)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.BookingDetails{}, domain.ErrBookingNotFound
		}
		return domain.BookingDetails{}, err
	}
	return goverter_postgres.ToDomainBookingDetails(row), nil
}
