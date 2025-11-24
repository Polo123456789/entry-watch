package app

import "time"

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
	CondoGetByID(id int64) (*Condominium, error)
	CondoCreate(condo *Condominium) (*Condominium, error)
	CondoUpdate(
		id int64,
		updateFn func(condo *Condominium) (*Condominium, error),
	) error
}
