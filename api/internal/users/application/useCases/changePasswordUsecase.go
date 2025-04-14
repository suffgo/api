package usecases

import (
	"fmt"
	sv "suffgo/internal/shared/domain/valueObjects"
	d "suffgo/internal/users/domain"
	v "suffgo/internal/users/domain/valueObjects"
)

type ChangePassword struct {
	repository d.UserRepository
}

func NewChangePasswordUsecase(repo d.UserRepository) *ChangePassword {
	return &ChangePassword{
		repository: repo,
	}
}

func (s *ChangePassword) Execute(id sv.ID, newPassword v.Password) error {

	user, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}

	hashedPassword, err := v.HashPassword(newPassword.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Obtener el ID del usuario existente
	userId := user.ID()

	// Crear usuario actualizado
	updateUser := d.NewUser(
		&userId,
		user.FullName(),
		user.Username(),
		user.Dni(),
		user.Email(),
		*hashedPassword,
		user.Image(),
	)

	// Actualizar en la base de datos
	_, err = s.repository.Update(*updateUser)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
