package usecases

import (
	"errors"
	"suffgo/internal/users/domain"
	valueobjects "suffgo/internal/users/domain/valueObjects"
)

type LoginUsecase struct {
	repository domain.UserRepository
}

// Este caso de uso es el encargado de validar el inicio de sesion del usuario
func NewLoginUsecase(repo domain.UserRepository) *LoginUsecase {
	return &LoginUsecase{
		repository: repo,
	}
}

func (s *LoginUsecase) Execute(
	username valueobjects.UserName,
	password valueobjects.Password,
) (*domain.User, error) {
	user, err := s.repository.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("credenciales invalidas")
	}

	if password.Password != user.Password().Password {
		return nil, errors.New("credenciales invalidas")
	}

	return user, nil
}
