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

func (m *mockQuerier) SearchHotelsByName(ctx context.Context, name string) ([]sqlc.Hotel, error) {
	args := m.Called(ctx, name)
	if rows, ok := args.Get(0).([]sqlc.Hotel); ok {
		return rows, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockQuerier) CreateRoom(ctx context.Context, arg sqlc.CreateRoomParams) (sqlc.Room, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(sqlc.Room); ok {
		return row, args.Error(1)
	}
	return sqlc.Room{}, args.Error(1)
}

func (m *mockQuerier) GetRoomByID(ctx context.Context, id uuid.UUID) (sqlc.Room, error) {
	args := m.Called(ctx, id)
	if row, ok := args.Get(0).(sqlc.Room); ok {
		return row, args.Error(1)
	}
	return sqlc.Room{}, args.Error(1)
}

func (m *mockQuerier) GetRoomForUpdate(ctx context.Context, id uuid.UUID) (sqlc.Room, error) {
	args := m.Called(ctx, id)
	if row, ok := args.Get(0).(sqlc.Room); ok {
		return row, args.Error(1)
	}
	return sqlc.Room{}, args.Error(1)
}

func (m *mockQuerier) GetRoomsByHotelID(ctx context.Context, hotelID uuid.UUID) ([]sqlc.Room, error) {
	args := m.Called(ctx, hotelID)
	if rows, ok := args.Get(0).([]sqlc.Room); ok {
		return rows, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockQuerier) TruncateAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockQuerier) FindAvailableRooms(
	ctx context.Context,
	arg sqlc.FindAvailableRoomsParams,
) ([]sqlc.FindAvailableRoomsRow, error) {
	args := m.Called(ctx, arg)
	if rows, ok := args.Get(0).([]sqlc.FindAvailableRoomsRow); ok {
		return rows, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockQuerier) CreateBooking(ctx context.Context, arg sqlc.CreateBookingParams) (sqlc.Booking, error) {
	args := m.Called(ctx, arg)
	if row, ok := args.Get(0).(sqlc.Booking); ok {
		return row, args.Error(1)
	}
	return sqlc.Booking{}, args.Error(1)
}

func (m *mockQuerier) CheckRoomBookingConflict(
	ctx context.Context,
	arg sqlc.CheckRoomBookingConflictParams,
) (int64, error) {
	args := m.Called(ctx, arg)
	val, _ := args.Get(0).(int64)
	return val, args.Error(1)
}

func (m *mockQuerier) GetBookingByReference(ctx context.Context, reference string) (sqlc.Booking, error) {
	args := m.Called(ctx, reference)
	if row, ok := args.Get(0).(sqlc.Booking); ok {
		return row, args.Error(1)
	}
	return sqlc.Booking{}, args.Error(1)
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
func TestHotelRepository_SearchByName_Success(t *testing.T) {
	q := &mockQuerier{}
	repo := NewHotelRepositoryWithQuerier(q)
	ctx := t.Context()

	hotelID := uuid.New()
	fixedTime := time.Date(2026, 10, 6, 12, 0, 0, 0, time.UTC)
	rows := []sqlc.Hotel{
		{
			ID:        hotelID,
			Name:      testHotelName,
			Address:   testHotelAddress,
			City:      testHotelCity,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		},
	}

	q.On("SearchHotelsByName", ctx, "Budapest").Return(rows, nil).Once()

	res, err := repo.SearchByName(ctx, "Budapest")
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, hotelID, res[0].ID)
	assert.Equal(t, testHotelName, res[0].Name)
	assert.Equal(t, testHotelAddress, res[0].Address)
	assert.Equal(t, testHotelCity, res[0].City)
	assert.True(t, fixedTime.Equal(res[0].CreatedAt))
	assert.True(t, fixedTime.Equal(res[0].UpdatedAt))

	q.AssertExpectations(t)
}

func TestHotelRepository_SearchByName_Error(t *testing.T) {
	q := &mockQuerier{}
	repo := NewHotelRepositoryWithQuerier(q)
	ctx := t.Context()

	expectedErr := errors.New("db error")
	q.On("SearchHotelsByName", ctx, "Budapest").Return(nil, expectedErr).Once()

	res, err := repo.SearchByName(ctx, "Budapest")
	require.ErrorIs(t, err, expectedErr)
	assert.Nil(t, res)

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
