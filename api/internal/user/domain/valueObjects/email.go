package valueobjects

import "errors"

type (
	Email struct {
		Email string
	}
)

func NewEmail(email string) (*Email, error) {
	if email == "" {
		return nil, errors.New("invalid email")
	}

	return &Email{
		Email: email,
	}, nil
}
