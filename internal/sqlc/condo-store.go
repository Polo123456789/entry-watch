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
	created, err := s.Queries.CondoCreate(ctx, CondoCreateParams{
		Name:      condo.Name,
		Address:   condo.Address,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: sql.NullInt64{Int64: condo.CreatedBy, Valid: condo.CreatedBy != 0},
		UpdatedBy: sql.NullInt64{Int64: condo.UpdatedBy, Valid: condo.UpdatedBy != 0},
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
	_, err = s.Queries.CondoUpdate(ctx, CondoUpdateParams{
		ID:        id,
		Name:      updated.Name,
		Address:   updated.Address,
		UpdatedAt: time.Now().Unix(),
		UpdatedBy: sql.NullInt64{Int64: updated.UpdatedBy, Valid: updated.UpdatedBy != 0},
	})
	return err
}

func (s *Store) CondoDelete(ctx context.Context, id int64) error {
	return s.Queries.CondoDelete(ctx, id)
}
