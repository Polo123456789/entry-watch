package sqlc

import (
	"context"
	"database/sql"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

// Store implements the entry.Store interface using SQLC queries.
type Store struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new Store with the given database connection.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// VisitGetByID retrieves a visit by its ID.
func (s *Store) VisitGetByID(ctx context.Context, id string) (*entry.Visit, error) {
	panic("Not implemented")
	return nil, nil
}

// VisitCreate creates a new visit.
func (s *Store) VisitCreate(ctx context.Context, visit *entry.Visit) (*entry.Visit, error) {
	panic("Not implemented")
	return visit, nil
}

// VisitUpdate updates an existing visit.
func (s *Store) VisitUpdate(
	ctx context.Context,
	id string,
	updateFn func(visit *entry.Visit) (*entry.Visit, error),
) error {
	panic("Not implemented")
	return nil
}

// CondoGetByID retrieves a condominium by its ID.
func (s *Store) CondoGetByID(ctx context.Context, id int64) (*entry.Condominium, error) {
	panic("Not implemented")
	return nil, nil
}

// CondoCreate creates a new condominium.
func (s *Store) CondoCreate(ctx context.Context, condo *entry.Condominium) (*entry.Condominium, error) {
	panic("Not implemented")
	return condo, nil
}

// CondoUpdate updates an existing condominium.
func (s *Store) CondoUpdate(
	ctx context.Context,
	id int64,
	updateFn func(condo *entry.Condominium) (*entry.Condominium, error),
) error {
	panic("Not implemented")
	return nil
}
