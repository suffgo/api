package usecases

import (
	"suffgo-backend-t/internal/user/entities"
	"suffgo-backend-t/internal/user/models"
)

type UserUsecase interface {
	UserDataRegister(in *models.AddUserData) error
	GetUserByID(id string) (*entities.UserSafeDto, error)
	DeleteUser(id string) error
	GetAll() ([]entities.UserSafeDto, error)
}
