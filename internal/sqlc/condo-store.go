package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

func (s *Store) CondoList(ctx context.Context) ([]*entry.Condominium, error) {
	condos, err := s.Queries.CondoList(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*entry.Condominium, len(condos))
	for i, c := range condos {
		result[i] = c.unmarshall()
	}
	return result, nil
}

func (s *Store) CondoGetByID(ctx context.Context, id int64) (*entry.Condominium, bool, error) {
	condo, err := s.Queries.CondoGetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return condo.unmarshall(), true, nil
}

func (s *Store) CondoCreate(ctx context.Context, condo *entry.Condominium) (*entry.Condominium, error) {
	now := time.Now().Unix()
	user := entry.UserFromCtx(ctx)
	var createdBy sql.NullInt64
	if user != nil {
		createdBy = sql.NullInt64{Int64: user.ID, Valid: true}
	}
	created, err := s.Queries.CondoCreate(ctx, CondoCreateParams{
		Name:      condo.Name,
		Address:   condo.Address,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	})
	if err != nil {
		return nil, err
	}
	return created.unmarshall(), nil
}

func (s *Store) CondoUpdate(
	ctx context.Context,
	id int64,
	updateFn func(condo *entry.Condominium) (*entry.Condominium, error),
) error {
	condo, ok, err := s.CondoGetByID(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return entry.NewUserSafeError("Condominio no encontrado")
	}
	updated, err := updateFn(condo)
	if err != nil {
		return err
	}
	user := entry.UserFromCtx(ctx)
	var updatedBy sql.NullInt64
	if user != nil {
		updatedBy = sql.NullInt64{Int64: user.ID, Valid: true}
	}
	_, err = s.Queries.CondoUpdate(ctx, CondoUpdateParams{
		ID:        id,
		Name:      updated.Name,
		Address:   updated.Address,
		UpdatedAt: time.Now().Unix(),
		UpdatedBy: updatedBy,
	})
	return err
}

func (s *Store) CondoDelete(ctx context.Context, id int64) error {
	return s.Queries.CondoDelete(ctx, id)
}
