package entry

import "context"

type UserRole string

const (
	RoleSuperAdmin UserRole = "superadmin"
	RoleAdmin      UserRole = "admin"
	RoleUser       UserRole = "user"
	RoleGuardian   UserRole = "guard"
)

type User struct {
	ID            int64
	CondominiumID int64
	Role          UserRole
	Enabled       bool
}

type userCtxKey struct{}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

func UserFromCtx(ctx context.Context) *User {
	user, _ := ctx.Value(userCtxKey{}).(*User)
	return user
}

type UnauthorizedError struct {
	msg string
}

func (e *UnauthorizedError) Error() string {
	return e.msg
}

type ForbiddenError struct {
	msg string
}

func (e *ForbiddenError) Error() string {
	return e.msg
}

func RequireRole(ctx context.Context, role UserRole) (*User, error) {
	user := UserFromCtx(ctx)
	if user == nil {
		return nil, &UnauthorizedError{msg: "user not authenticated"}
	}
	if !user.Enabled {
		return nil, &ForbiddenError{msg: "user is disabled"}
	}
	if user.Role != role && user.Role != RoleSuperAdmin {
		return nil, &ForbiddenError{msg: "insufficient permissions"}
	}
	return user, nil
}

func RequireRoleAndCondo(
	ctx context.Context, role UserRole, condoID int64,
) (*User, error) {
	user, err := RequireRole(ctx, role)
	if err != nil {
		return nil, err
	}
	if user.CondominiumID != condoID && user.Role != RoleSuperAdmin {
		return nil, &ForbiddenError{
			msg: "insufficient permissions for this condominium",
		}
	}
	return user, nil
}

// UserWithPassword extends User with the password hash for authentication.
type UserWithPassword struct {
	*User
	PasswordHash string
}

// UserSafeError represents an error that can be safely shown to users.
// These errors are typically caused by user input and don't expose system details.
type UserSafeError struct {
	Msg string
}

func (e *UserSafeError) Error() string {
	return e.Msg
}

// NewUserSafeError creates a new user-safe error with the given message.
func NewUserSafeError(msg string) *UserSafeError {
	return &UserSafeError{Msg: msg}
}
