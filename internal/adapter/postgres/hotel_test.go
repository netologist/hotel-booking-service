package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/netologist/hotel-booking/internal/domain"
	"github.com/netologist/hotel-booking/internal/gen/sqlc"
)

const (
	testHotelName    = "Grand Budapest Hotel"
	testHotelAddress = "100 Alpine Way"
	testHotelCity    = "Lutz"
)

type mockQuerier struct {
	mock.Mock
}

var _ sqlc.Querier = (*mockQuerier)(nil)

func (m *mockQuerier) CreateHotel(ctx context.Context, arg sqlc.CreateHotelParams) (sqlc.Hotel, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(sqlc.Hotel); ok {
		return row, args.Error(1)
	}
	return sqlc.Hotel{}, args.Error(1)
}

func TestNewHotelRepository(t *testing.T) {
	repo := NewHotelRepository(nil)
	require.NotNil(t, repo)
}

func TestNewHotelRepositoryWithQuerier(t *testing.T) {
	q := &mockQuerier{}
	repo := NewHotelRepositoryWithQuerier(q)
	require.NotNil(t, repo)
}

func TestHotelRepository_Create_Success(t *testing.T) {
	q := &mockQuerier{}
	repo := NewHotelRepositoryWithQuerier(q)
	ctx := t.Context()

	hotelID := uuid.New()
	fixedTime := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	input := domain.Hotel{
		ID:      hotelID,
		Name:    testHotelName,
		Address: testHotelAddress,
		City:    testHotelCity,
	}

	expectedArg := sqlc.CreateHotelParams{
		ID:      hotelID,
		Name:    testHotelName,
		Address: testHotelAddress,
		City:    testHotelCity,
	}

	createdRow := sqlc.Hotel{
		ID:        hotelID,
		Name:      testHotelName,
		Address:   testHotelAddress,
		City:      testHotelCity,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}
	q.On("CreateHotel", ctx, expectedArg).Return(createdRow, nil).Once()

	res, err := repo.Create(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, hotelID, res.ID)
	assert.Equal(t, testHotelName, res.Name)
	assert.Equal(t, testHotelAddress, res.Address)
	assert.Equal(t, testHotelCity, res.City)
	assert.True(t, fixedTime.Equal(res.CreatedAt))
	assert.True(t, fixedTime.Equal(res.UpdatedAt))

	q.AssertExpectations(t)
}

func TestHotelRepository_Create_Error(t *testing.T) {
	q := &mockQuerier{}
	repo := NewHotelRepositoryWithQuerier(q)
	ctx := t.Context()
	dbErr := errors.New("unique constraint violation")

	hotelID := uuid.New()
	input := domain.Hotel{
		ID:      hotelID,
		Name:    testHotelName,
		Address: testHotelAddress,
		City:    testHotelCity,
	}

	expectedArg := sqlc.CreateHotelParams{
		ID:      hotelID,
		Name:    testHotelName,
		Address: testHotelAddress,
		City:    testHotelCity,
	}
	q.On("CreateHotel", ctx, expectedArg).Return(sqlc.Hotel{}, dbErr).Once()

	res, err := repo.Create(ctx, input)
	require.ErrorIs(t, err, dbErr)
	assert.Zero(t, res)

	q.AssertExpectations(t)
}

func TestHotelRepository_Create_UniqueViolation(t *testing.T) {
	q := &mockQuerier{}
	repo := NewHotelRepositoryWithQuerier(q)
	ctx := t.Context()
	pgErr := &pgconn.PgError{
		Code:   "23505",
		Detail: "Key (name)=(Grand Budapest Hotel) already exists.",
	}

	hotelID := uuid.New()
	input := domain.Hotel{
		ID:      hotelID,
		Name:    testHotelName,
		Address: testHotelAddress,
		City:    testHotelCity,
	}

	expectedArg := sqlc.CreateHotelParams{
		ID:      hotelID,
		Name:    testHotelName,
		Address: testHotelAddress,
		City:    testHotelCity,
	}
	q.On("CreateHotel", ctx, expectedArg).Return(sqlc.Hotel{}, pgErr).Once()

	res, err := repo.Create(ctx, input)
	require.ErrorIs(t, err, domain.ErrHotelAlreadyExists)
	assert.Zero(t, res)

	q.AssertExpectations(t)
}

func TestToDomainHotel(t *testing.T) {
	hotelID := uuid.New()
	fixedTime := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	row := sqlc.Hotel{
		ID:        hotelID,
		Name:      testHotelName,
		Address:   testHotelAddress,
		City:      testHotelCity,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}
	got := toDomainHotel(row)
	require.NotNil(t, got)

	assert.Equal(t, hotelID, got.ID)
	assert.Equal(t, testHotelName, got.Name)
	assert.Equal(t, testHotelAddress, got.Address)
	assert.Equal(t, testHotelCity, got.City)
	assert.True(t, fixedTime.Equal(got.CreatedAt))
	assert.True(t, fixedTime.Equal(got.UpdatedAt))
}
