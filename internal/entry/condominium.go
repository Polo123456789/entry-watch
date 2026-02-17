package entry

import (
	"context"
	"time"
)

type Condominium struct {
	ID        int64
	Name      string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy int64
	UpdatedBy int64
}

type CondominiumStore interface {
	CondoList(ctx context.Context) ([]*Condominium, error)
	CondoGetByID(ctx context.Context, id int64) (*Condominium, error)
	CondoCreate(ctx context.Context, condo *Condominium) (*Condominium, error)
	CondoUpdate(
		ctx context.Context,
		id int64,
		updateFn func(condo *Condominium) (*Condominium, error),
	) error
	CondoDelete(ctx context.Context, id int64) error
}
