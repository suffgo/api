package valueobjects

import (
	"errors"
	"strings"
)

type (
	Email struct {
		Email string
	}
)

func NewEmail(email string) (*Email, error) {
	if email == "" {
		return nil, errors.New("invalid email")
	}

	if strings.Contains(email, "@") == false {
		return nil, errors.New("invalid email")
	}

	if strings.Contains(email, ".") == false {
		return nil, errors.New("invalid email")
	}

	if strings.Count(email, "@") != 1 {
		return nil, errors.New("invalid email")
	}

	parts := strings.Split(email, "@")
	local := parts[0]
	dominio := parts[1]

	if strings.HasPrefix(dominio, ".") || strings.HasSuffix(dominio, ".") {
		return nil, errors.New("invalid email")
	}

	if len(local) == 0 {
		return nil, errors.New("invalid email")
	}


	return &Email{
		Email: email,
	}, nil
}
