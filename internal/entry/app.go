package entry

import (
	"context"
	"log/slog"
)

type App struct {
	Config Config
	store  Store
	logger *slog.Logger
}

func NewApp(logger *slog.Logger, store Store) *App {
	return &App{
		store:  store,
		logger: logger,
		Config: Config{},
	}
}

type Store interface {
	CondominiumStore
	VisitStore
}

type Config struct{}

type Valid interface {
	Valid() error
}

func (a *App) CondoList(ctx context.Context) ([]*Condominium, error) {
	return a.store.CondoList(ctx)
}

func (a *App) CondoGetByID(ctx context.Context, id int64) (*Condominium, error) {
	return a.store.CondoGetByID(ctx, id)
}

func (a *App) CondoCreate(ctx context.Context, name, address string, createdBy int64) (*Condominium, error) {
	condo := &Condominium{
		Name:      name,
		Address:   address,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
	}
	if err := condo.Valid(); err != nil {
		return nil, err
	}
	return a.store.CondoCreate(ctx, condo)
}

func (a *App) CondoUpdate(ctx context.Context, id int64, name, address string, updatedBy int64) error {
	return a.store.CondoUpdate(ctx, id, func(condo *Condominium) (*Condominium, error) {
		condo.Name = name
		condo.Address = address
		condo.UpdatedBy = updatedBy
		if err := condo.Valid(); err != nil {
			return nil, err
		}
		return condo, nil
	})
}

func (a *App) CondoDelete(ctx context.Context, id int64) error {
	return a.store.CondoDelete(ctx, id)
}
