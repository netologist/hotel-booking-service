package postgres

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAdminRepository(t *testing.T) {
	repo := NewAdminRepository(nil)
	require.NotNil(t, repo)
}

func TestAdminRepository_TruncateAll_Success(t *testing.T) {
	q := &mockQuerier{}
	repo := NewAdminRepositoryWithQuerier(q)
	ctx := t.Context()

	q.On("TruncateAll", ctx).Return(nil).Once()

	err := repo.TruncateAll(ctx)
	require.NoError(t, err)
	q.AssertExpectations(t)
}

func TestAdminRepository_TruncateAll_Error(t *testing.T) {
	q := &mockQuerier{}
	repo := NewAdminRepositoryWithQuerier(q)
	ctx := t.Context()

	dbErr := errors.New("db error")
	q.On("TruncateAll", ctx).Return(dbErr).Once()

	err := repo.TruncateAll(ctx)
	require.ErrorIs(t, err, dbErr)
	q.AssertExpectations(t)
}
