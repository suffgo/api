package errors

type notInWhitelist string

const ErrNotWhitelist notInWhitelist = "the current user is not in whtelist."

func (r notInWhitelist) Error() string {
	return string(r)
}