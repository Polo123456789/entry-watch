package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

type Store struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) VisitGetByID(ctx context.Context, id string) (*entry.Visit, error) {
	panic("Not implemented")
}

func (s *Store) VisitCreate(ctx context.Context, visit *entry.Visit) (*entry.Visit, error) {
	panic("Not implemented")
}

func (s *Store) VisitUpdate(
	ctx context.Context,
	id string,
	updateFn func(visit *entry.Visit) (*entry.Visit, error),
) error {
	panic("Not implemented")
}

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

func (s *Store) CondoGetByID(ctx context.Context, id int64) (*entry.Condominium, error) {
	condo, err := s.Queries.CondoGetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return condo.unmarshall(), nil
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
	condo, err := s.Queries.CondoGetByID(ctx, id)
	if err != nil {
		return err
	}
	updated, err := updateFn(condo.unmarshall())
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
