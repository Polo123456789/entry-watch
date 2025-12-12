package entry

import (
	"context"
	"time"
)

type Visit struct {
	ID            string
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
	VisitGetByID(ctx context.Context, id string) (*Visit, error)
	VisitCreate(ctx context.Context, visit *Visit) (*Visit, error)
	VisitUpdate(
		ctx context.Context,
		id string,
		updateFn func(visit *Visit) (*Visit, error),
	) error
}
