package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/netologist/hotel-booking/internal/domain"
	"github.com/netologist/hotel-booking/internal/gen/sqlc"
)

type adminRepository struct {
	queries sqlc.Querier
}

var _ domain.AdminRepository = (*adminRepository)(nil)

func NewAdminRepository(pool *pgxpool.Pool) *adminRepository {
	return NewAdminRepositoryWithQuerier(sqlc.New(pool))
}

func NewAdminRepositoryWithQuerier(queries sqlc.Querier) *adminRepository {
	return &adminRepository{queries: queries}
}

func (a *adminRepository) TruncateAll(ctx context.Context) error {
	return a.queries.TruncateAll(ctx)
}
