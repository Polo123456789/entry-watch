package sqlc

import (
	"github.com/Polo123456789/entry-watch/internal/entry"
	"github.com/Polo123456789/entry-watch/internal/http/auth"
)

// unmarshall converts a SQLC User to an auth.User.
func (u User) unmarshall() *auth.User {
	return &auth.User{
		ID:            u.ID,
		CondominiumID: getInt64FromNullInt64(u.CondominiumID),
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		Phone:         getStringFromNullString(u.Phone),
		Role:          entry.UserRole(u.Role),
		Enabled:       u.Enabled,
		Hidden:        u.Hidden,
	}
}

// unmarshallSlice converts a slice of SQLC Users to a slice of auth.Users.
//
//nolint:unused // Available for future use when implementing user listing
func unmarshallSlice(users []User) []*auth.User {
	result := make([]*auth.User, len(users))
	for i, u := range users {
		result[i] = u.unmarshall()
	}
	return result
}
