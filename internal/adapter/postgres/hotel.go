package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/netologist/hotel-booking/internal/domain"
	goverter_postgres "github.com/netologist/hotel-booking/internal/gen/goverter/postgres"
	"github.com/netologist/hotel-booking/internal/gen/sqlc"
)

// HotelRepository implements domain.HotelRepository via SQLC.
type hotelRepository struct {
	queries sqlc.Querier
}

var _ domain.HotelRepository = (*hotelRepository)(nil)

// NewHotelRepository creates a new hotelRepository.
func NewHotelRepository(pool *pgxpool.Pool) *hotelRepository {
	return NewHotelRepositoryWithQuerier(sqlc.New(pool))
}

// NewHotelRepositoryWithQuerier creates a new hotelRepository with a custom Querier.
func NewHotelRepositoryWithQuerier(queries sqlc.Querier) *hotelRepository {
	return &hotelRepository{queries: queries}
}

func (h *hotelRepository) Create(ctx context.Context, hotel domain.Hotel) (domain.Hotel, error) {
	row, err := h.queries.CreateHotel(ctx, goverter_postgres.ToCreateHotelParams(hotel))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.Hotel{}, fmt.Errorf("%w: %w", domain.ErrHotelAlreadyExists, err)
		}
		return domain.Hotel{}, err
	}
	return toDomainHotel(row), nil
}

func (h *hotelRepository) SearchByName(ctx context.Context, name string) ([]domain.Hotel, error) {
	rows, err := h.queries.SearchHotelsByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return goverter_postgres.ToDomainHotels(rows), nil
}

// Model mappers
func toDomainHotel(row sqlc.Hotel) domain.Hotel {
	return goverter_postgres.ToDomainHotel(row)
}
