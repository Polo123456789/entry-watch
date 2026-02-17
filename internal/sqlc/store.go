package sqlc

import (
	"context"
	"database/sql"

	"github.com/Polo123456789/entry-watch/internal/entry"
)

type Store struct {
	db *sql.DB
	*Queries
}

var _ entry.Store = (*Store)(nil)

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) VisitGetByID(ctx context.Context, id string) (*entry.Visit, bool, error) {
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
