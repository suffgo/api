package usecases

import (
	"suffgo/internal/user/entities"
	"suffgo/internal/user/models"
)

type UserUsecase interface {
	UserDataRegister(in *models.AddUserData) error
	GetUserByID(id string) (*entities.UserSafeDto, error)
	DeleteUser(id string) error
	GetAll() ([]entities.UserSafeDto, error)
}
