package mappers

import (
	"suffgo/internal/user/domain"
	v "suffgo/internal/user/domain/valueObjects"
	m "suffgo/internal/user/infrastructure/models"
)

func DomainToModel(user *domain.User) *m.User {
	return &m.User{
		ID:       user.ID().Id, // Convierte UserID a uint
		Dni:      user.Dni().Dni,
		Username: user.Username().Username,
		Password: user.Password().Password,
		Name:     user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Email:    user.Email().Email,
	}
}

func modelToDomain(userModel *m.User) (*domain.User, error) {
	id := v.NewUserID(userModel.ID)
	name := v.NewUserFullName(userModel.Name, userModel.Lastname)
	username := v.NewUserUserName(userModel.Username)
	dni := v.NewUserDni(userModel.Dni)
	email := v.NewUserEmail(userModel.Email)
	password := v.NewUserPassword(userModel.Password)

	return domain.NewUser(id, *name, *username, *dni, *email, *password), nil
}
