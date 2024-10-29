package usecases

import (
	"errors"
	"suffgo/internal/user/domain"
	valueobjects "suffgo/internal/user/domain/valueObjects"
)

type LoginUsecase struct {
	repository domain.UserRepository
}

func NewLoginUsecase(repo domain.UserRepository) *LoginUsecase {
	return &LoginUsecase{
		repository: repo,
	}
}

func (s *LoginUsecase) Execute(username valueobjects.UserName, password valueobjects.Password) error {

	user, err := s.repository.GetByUsername(username)

	if err != nil {
		return err
	}

	if password.Password == user.Password().Password {
		return nil
	}
	
	return errors.New("Credenciales invalidas")
}