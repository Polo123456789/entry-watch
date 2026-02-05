package entry

// UserSafeError represents an error that can be safely shown to users.
// These errors are typically caused by user input and don't expose system
// details.
type UserSafeError struct {
	msg string
}

func (e UserSafeError) Error() string {
	return e.msg
}

func NewUserSafeError(msg string) UserSafeError {
	return UserSafeError{msg: msg}
}
