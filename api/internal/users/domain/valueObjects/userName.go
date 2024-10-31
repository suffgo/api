package valueobjects

import "errors"

type (
	UserName struct {
		Username string
	}
)

func NewUserName(username string) (*UserName, error) {
	if username == "" {
		return nil, errors.New("invalid username")
	}

	return &UserName{
		Username: username,
	}, nil
}
