package repositories

import "suffgo/internal/user/entities"

type UserRepository interface {
	InsertUserData(in *entities.UserDto) error
	GetUserByID(id int) (*entities.UserSafeDto, error)
	DeleteUser(id int) error
	FetchAll() ([]entities.User, error)
}
