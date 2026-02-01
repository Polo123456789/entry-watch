package auth

import (
	"context"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/sqlc"
)

// sqlcUserStore adapts sqlc.UserStore to auth.UserStore interface.
type sqlcUserStore struct {
	store *sqlc.UserStore
}

// NewSQLCUserStore creates a UserStore from a sqlc.UserStore.
func NewSQLCUserStore(store *sqlc.UserStore) UserStore {
	return &sqlcUserStore{store: store}
}

func (s *sqlcUserStore) GetByEmail(ctx context.Context, email string) (UserWithPassword, bool, error) {
	user, passwordHash, ok, err := s.store.GetByEmail(ctx, email)
	if err != nil || !ok {
		return UserWithPassword{}, ok, err
	}

	return UserWithPassword{
		User:         user,
		PasswordHash: passwordHash,
	}, true, nil
}

func (s *sqlcUserStore) GetByID(ctx context.Context, id int64) (*entry.User, bool, error) {
	return s.store.GetByID(ctx, id)
}

func (s *sqlcUserStore) CreateUser(ctx context.Context, email, firstName, lastName string, user *entry.User, passwordHash string) error {
	return s.store.CreateUser(ctx, email, firstName, lastName, user, passwordHash)
}

func (s *sqlcUserStore) CountSuperAdmins(ctx context.Context) (int64, error) {
	return s.store.CountSuperAdmins(ctx)
}
