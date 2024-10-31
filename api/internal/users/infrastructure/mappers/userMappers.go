package mappers

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	"suffgo/internal/users/domain"
	v "suffgo/internal/users/domain/valueObjects"
	m "suffgo/internal/users/infrastructure/models"
)

func DomainToModel(user *domain.User) *m.Users {
	return &m.Users{
		ID:       user.ID().Id, // Convierte ID a uint
		Dni:      user.Dni().Dni,
		Username: user.Username().Username,
		Password: user.Password().Password,
		Name:     user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Email:    user.Email().Email,
	}
}

func ModelToDomain(userModel *m.Users) (*domain.User, error) {
	id, err := sv.NewID(userModel.ID)
	if err != nil {
		return nil, err
	}
	name, err := v.NewFullName(userModel.Name, userModel.Lastname)
	if err != nil {
		return nil, err
	}
	username, err := v.NewUserName(userModel.Username)
	if err != nil {
		return nil, err
	}
	dni, err := v.NewDni(userModel.Dni)
	if err != nil {
		return nil, err
	}
	email, err := v.NewEmail(userModel.Email)
	if err != nil {
		return nil, err
	}
	password, err := v.NewPassword(userModel.Password)
	if err != nil {
		return nil, err
	}

	return domain.NewUser(id, *name, *username, *dni, *email, *password), nil
}
