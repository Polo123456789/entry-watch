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

func (c *Condominium) Valid() error {
	if len(c.Name) == 0 || len(c.Name) > 200 {
		return NewUserSafeError("El nombre debe tener entre 1 y 200 caracteres")
	}
	if len(c.Address) == 0 || len(c.Address) > 500 {
		return NewUserSafeError("La dirección debe tener entre 1 y 500 caracteres")
	}
	return nil
}

type CondominiumStore interface {
	CondoList(ctx context.Context) ([]*Condominium, error)
	CondoGetByID(ctx context.Context, id int64) (*Condominium, bool, error)
	CondoCreate(ctx context.Context, condo *Condominium) (*Condominium, error)
	CondoUpdate(
		ctx context.Context,
		id int64,
		updateFn func(condo *Condominium) (*Condominium, error),
	) error
	CondoDelete(ctx context.Context, id int64) error
}

func (a *App) CondoList(ctx context.Context) ([]*Condominium, error) {
	return a.store.CondoList(ctx)
}

func (a *App) CondoGetByID(ctx context.Context, id int64) (*Condominium, bool, error) {
	return a.store.CondoGetByID(ctx, id)
}

func (a *App) CondoCreate(ctx context.Context, name, address string) (*Condominium, error) {
	condo := &Condominium{
		Name:    name,
		Address: address,
	}
	if err := condo.Valid(); err != nil {
		return nil, err
	}
	return a.store.CondoCreate(ctx, condo)
}

func (a *App) CondoUpdate(ctx context.Context, id int64, name, address string) error {
	return a.store.CondoUpdate(
		ctx, id,
		func(condo *Condominium) (*Condominium, error) {
			condo.Name = name
			condo.Address = address
			if err := condo.Valid(); err != nil {
				return nil, err
			}
			return condo, nil
		},
	)
}

func (a *App) CondoDelete(ctx context.Context, id int64) error {
	return a.store.CondoDelete(ctx, id)
}
