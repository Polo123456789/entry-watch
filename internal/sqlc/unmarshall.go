package sqlc

import (
	"database/sql"

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
