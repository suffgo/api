package usecases

import (
	"fmt"
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

func (s *ChangePassword) Execute(email v.Email, newPassword v.Password) error {
	// Buscar usuario por email
	user, err := s.repository.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Generar hash de la nueva contraseña
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
	_, err = s.repository.Update(*updateUser) // Manejar ambos valores de retorno
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
