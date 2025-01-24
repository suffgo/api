package errors

type roomAlreadyInWhitelist string

const ErrAlreadyInWhitelist roomAlreadyInWhitelist = "user already in whitelist."

func (r roomAlreadyInWhitelist) Error() string {
	return string(r)
}
