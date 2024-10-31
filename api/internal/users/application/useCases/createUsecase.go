package usecases

import (
	"errors"
	"suffgo/internal/users/domain"
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

	existingUser, err = s.repository.GetByDni(user.Dni())
	if err != nil {
		return err
	}

	if existingUser != nil {
		// Usuario con mismo dni ya existe
		return errors.New("user already exists with this dni")
	}

	existingUser, err = s.repository.GetByUsername(user.Username())
	if err != nil {
		return err
	}

	if existingUser != nil {
		// Usuario con mismo username  ya existe
		return errors.New("user already exists with this username")
	}

	// Save the user to the repository
	err = s.repository.Save(user)
	if err != nil {
		return err
	}

	return nil
}
