package valueobjects

import (
	"errors"
	"strings"
)

type (
	UserName struct {
		Username string
	}
)

func NewUserName(username string) (*UserName, error) {
	if username == "" {
		return nil, errors.New("invalid username")
	}

	usernameMin := strings.ToLower(username)

	return &UserName{
		Username: usernameMin,
	}, nil
}
