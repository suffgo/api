package usecases

import (
	"errors"
	sv "suffgo/internal/shared/domain/valueObjects"
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

func (s *LoginUsecase) Execute(username valueobjects.UserName, password valueobjects.Password) (*sv.ID, error) {

	user, err := s.repository.GetByUsername(username)

	if err != nil {
		return nil, errors.New("Credenciales invalidas")
	}

	if password.Password != user.Password().Password {
		return nil, errors.New("Credenciales invalidas")
	}

	id := user.ID()

	return &id, nil
}
