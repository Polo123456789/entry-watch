package sqlc

import (
	"database/sql"
	"time"

	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
)

// validNullInt64 safely extracts int64 from sql.NullInt64.
func validNullInt64(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}

// validNullString safely extracts string from sql.NullString.
func validNullString(n sql.NullString) string {
	if n.Valid {
		return n.String
	}
	return ""
}

func (u User) unmarshall() *auth.User {
	return &auth.User{
		ID:            u.ID,
		CondominiumID: validNullInt64(u.CondominiumID),
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		Phone:         validNullString(u.Phone),
		Role:          entry.UserRole(u.Role),
		Enabled:       u.Enabled,
		Hidden:        u.Hidden,
	}
}

func (c Condominium) unmarshall() *entry.Condominium {
	return &entry.Condominium{
		ID:        c.ID,
		Name:      c.Name,
		Address:   c.Address,
		CreatedAt: time.Unix(c.CreatedAt, 0),
		UpdatedAt: time.Unix(c.UpdatedAt, 0),
		CreatedBy: validNullInt64(c.CreatedBy),
		UpdatedBy: validNullInt64(c.UpdatedBy),
	}
}
