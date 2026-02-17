package superadmin

import (
	"context"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
	"github.com/Polo123456789/entry-watch/internal/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type StoreAdapter struct {
	userStore *sqlc.UserStore
	appStore  entry.Store
}

func NewStoreAdapter(userStore *sqlc.UserStore, appStore entry.Store) *StoreAdapter {
	return &StoreAdapter{
		userStore: userStore,
		appStore:  appStore,
	}
}

func (s *StoreAdapter) UserListByRole(ctx context.Context, role entry.UserRole) ([]*entry.AdminUser, error) {
	users, err := s.userStore.UserListByRole(ctx, role)
	if err != nil {
		return nil, err
	}
	result := make([]*entry.AdminUser, len(users))
	for i, u := range users {
		result[i] = &entry.AdminUser{
			ID:            u.ID,
			CondominiumID: u.CondominiumID,
			FirstName:     u.FirstName,
			LastName:      u.LastName,
			Email:         u.Email,
			Phone:         u.Phone,
			Enabled:       u.Enabled,
			CondoName:     u.CondoName,
		}
	}
	return result, nil
}

func (s *StoreAdapter) UserGetByID(ctx context.Context, id int64) (*entry.AdminUser, bool, error) {
	user, ok, err := s.userStore.GetByID(ctx, id)
	if err != nil || !ok {
		return nil, false, err
	}
	return &entry.AdminUser{
		ID:            user.ID,
		CondominiumID: user.CondominiumID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Phone:         user.Phone,
		Enabled:       user.Enabled,
	}, true, nil
}

func (s *StoreAdapter) UserCreate(ctx context.Context, user *entry.AdminUser, password string) (*entry.AdminUser, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	authUser := &auth.User{
		CondominiumID: user.CondominiumID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Phone:         user.Phone,
		Role:          entry.RoleAdmin,
		Enabled:       user.Enabled,
		Hidden:        false,
	}

	created, err := s.userStore.CreateUser(ctx, authUser, string(hash))
	if err != nil {
		return nil, err
	}

	return &entry.AdminUser{
		ID:            created.ID,
		CondominiumID: created.CondominiumID,
		FirstName:     created.FirstName,
		LastName:      created.LastName,
		Email:         created.Email,
		Phone:         created.Phone,
		Enabled:       created.Enabled,
	}, nil
}

func (s *StoreAdapter) UserUpdate(ctx context.Context, id int64, user *entry.AdminUser, updatedBy int64) (*entry.AdminUser, error) {
	authUser := &auth.User{
		CondominiumID: user.CondominiumID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Phone:         user.Phone,
		Enabled:       user.Enabled,
	}

	updated, err := s.userStore.UserUpdate(ctx, id, authUser, updatedBy)
	if err != nil {
		return nil, err
	}

	return &entry.AdminUser{
		ID:            updated.ID,
		CondominiumID: updated.CondominiumID,
		FirstName:     updated.FirstName,
		LastName:      updated.LastName,
		Email:         updated.Email,
		Phone:         updated.Phone,
		Enabled:       updated.Enabled,
	}, nil
}

func (s *StoreAdapter) UserDelete(ctx context.Context, id int64) error {
	return s.userStore.UserDelete(ctx, id)
}

func (s *StoreAdapter) CondoList(ctx context.Context) ([]*entry.Condominium, error) {
	return s.appStore.CondoList(ctx)
}
