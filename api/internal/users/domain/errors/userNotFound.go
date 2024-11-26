package errors

// este error se usa cuando se solicita un usuario que no existe
type userNotFoundConst string

const ErrUserNotFound userNotFoundConst = "user not found."

func (u userNotFoundConst) Error() string {
	return string(u)
}
