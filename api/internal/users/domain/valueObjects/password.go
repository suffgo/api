package valueobjects

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

type (
	Password struct {
		Password string
	}
)

func NewPassword(password string) (*Password, error) {
	err := validation.Validate(password,
		validation.Required.Error("la contraseña es obligatoria"),
		validation.Length(8, 0).Error("la contraseña debe tener al menos 8 caracteres"),
		validation.Match(regexp.MustCompile(`[A-Z]`)).Error("la contraseña debe contener al menos una letra mayúscula"),
		validation.Match(regexp.MustCompile(`[a-z]`)).Error("la contraseña debe contener al menos una letra minúscula"),
		validation.Match(regexp.MustCompile(`[0-9]`)).Error("la contraseña debe contener al menos un número"),
	)
	if err != nil {
		return nil, err
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
