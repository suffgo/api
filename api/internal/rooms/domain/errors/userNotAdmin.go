package errors

type userNotAdminConst string

const ErrUserNotAdmin userNotAdminConst = "insufficient permissions."

func (u userNotAdminConst) Error() string {
	return string(u)
}
