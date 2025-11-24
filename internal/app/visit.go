package app

import (
	"context"
	"time"
)

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
	VisitGetByCode(ctx context.Context, code string) (*Visit, error)
	VisitCreate(ctx context.Context, visit *Visit) (*Visit, error)
	VisitUpdate(
		ctx context.Context,
		id int64,
		updateFn func(visit *Visit) (*Visit, error),
	) error
}
