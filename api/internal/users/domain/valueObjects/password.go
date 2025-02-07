package valueobjects

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type (
	Password struct {
		Password string
	}
)

func NewPassword(password string) (*Password, error) {
	if password == "" {
		return nil, errors.New("invalid password")
	}

	return &Password{
		Password: password,
	}, nil
}

func HashPassword(password string) (*Password, error) {
	if password == "" {
		return nil, errors.New("invalid password")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error generating password hash")
	}

	return &Password{
		Password: string(hashed),
	}, nil
}

func (p Password) Validate(password Password) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(password.Value()))
	return err == nil
}

func (p *Password) Value() string {
	return p.Password
}
