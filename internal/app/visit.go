package app

import "time"

type Visit struct {
	ID            int64
	Code          string
	CondominiumID int64
	UserID        int64
	VisitorName   string
	MaxUses       int64
	Uses          int64
	ValidFrom     time.Time
	ValidTo       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type VisitStore interface {
	VisitGetByCode(code string) (*Visit, error)
	VisitCreate(visit *Visit) (*Visit, error)
	VisitUpdate(
		id int64,
		updateFn func(visit *Visit) (*Visit, error),
	) error
}
