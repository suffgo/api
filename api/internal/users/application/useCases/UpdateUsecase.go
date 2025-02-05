package usecases

import (
	"errors"
	"suffgo/internal/users/domain"
	e "suffgo/internal/users/domain/errors"
)

type UpdateUsecase struct {
	repository domain.UserRepository
}

func NewUpdateUsecase(repository domain.UserRepository) *UpdateUsecase {
	return &UpdateUsecase{
		repository: repository,
	}
}

func (s *UpdateUsecase) Execute(user *domain.User) (*domain.User, error) {
	// Buscar usuario por ID
	existingUser, err := s.repository.GetByID(user.ID())
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, e.ErrUserNotFound
	}

	// Verificar si ya existe otro usuario con el mismo email
	existingUserByEmail, err := s.repository.GetByEmail(user.Email())
	if err != nil {
		return nil, err
	}
	if existingUserByEmail != nil && existingUserByEmail.ID() != user.ID() {
		return nil, errors.New("user already exists with this email")
	}

	// Verificar si ya existe otro usuario con el mismo DNI
	existingUserByDni, err := s.repository.GetByDni(user.Dni())
	if err != nil {
		return nil, err
	}
	if existingUserByDni != nil && existingUserByDni.ID() != user.ID() {
		return nil, errors.New("user already exists with this DNI")
	}

	// Verificar si ya existe otro usuario con el mismo username
	existingUserByUsername, err := s.repository.GetByUsername(user.Username())
	if err != nil {
		return nil, err
	}
	if existingUserByUsername != nil && existingUserByUsername.ID() != user.ID() {
		return nil, errors.New("user already exists with this username")
	}

	// Actualizar el usuario en la base de datos
	updatedUser, err := s.repository.Update(*user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
