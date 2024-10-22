package usecases

import (
	"errors"
	"suffgo/internal/user/domain"
)

type (
	CreateUsecase struct {
		repository domain.UserRepository
	}
)

func NewCreateUsecase(repository domain.UserRepository) *CreateUsecase {
	return &CreateUsecase{
		repository: repository,
	}
}

func (s *CreateUsecase) Execute(user domain.User) error {
	// Business Logic: Check if the user already exists
	existingUser, err := s.repository.GetByEmail(user.Email())
	if err != nil {
		return err
	}

	if existingUser != nil {
		// User with the same email already exists
		return errors.New("user already exists with this email")
	}

	// Save the user to the repository
	err = s.repository.Save(user)
	if err != nil {
		return err
	}

	return nil
}
