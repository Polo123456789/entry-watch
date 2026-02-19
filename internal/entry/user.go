package entry

import "context"

type UserRole string

const (
	RoleSuperAdmin UserRole = "superadmin"
	RoleAdmin      UserRole = "admin"
	RoleUser       UserRole = "user"
	RoleGuardian   UserRole = "guard"
)

// User represents a user in the domain layer.
// Contains only the information needed for domain-level authorization.
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

func NewUnauthorizedError(msg string) *UnauthorizedError {
	return &UnauthorizedError{msg: msg}
}

type ForbiddenError struct {
	msg string
}

func (e *ForbiddenError) Error() string {
	return e.msg
}

func NewForbiddenError(msg string) *ForbiddenError {
	return &ForbiddenError{msg: msg}
}

func RequireRole(ctx context.Context, role UserRole) (*User, error) {
	user := UserFromCtx(ctx)
	if user == nil {
		return nil, NewUnauthorizedError("user not authenticated")
	}
	if !user.Enabled {
		return nil, NewForbiddenError("user is disabled")
	}
	if user.Role != role && user.Role != RoleSuperAdmin {
		return nil, NewForbiddenError("insufficient permissions")
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
		return nil, NewForbiddenError("insufficient permissions for this condominium")
	}
	return user, nil
}
