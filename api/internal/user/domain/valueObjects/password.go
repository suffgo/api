package valueobjects

import "errors"

type (
	Password struct {
		Password string
	}
)

func NewPassword(password string) (*Password, error) {
	
	if password == "" {
		return nil, errors.New("Invalid password")
	}

	return &Password{
		Password: password,
	}, nil
}