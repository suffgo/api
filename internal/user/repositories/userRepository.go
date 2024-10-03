package repositories

import "suffgo-backend-t/internal/user/entities"

type UserRepository interface {
	InsertUserData(in *entities.UserDto) error
	GetUserByID(id int) (*entities.UserSafeDto, error)
	DeleteUser(id int) error
	FetchAll() ([]entities.User, error)
}
