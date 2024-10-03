package usecases

import (
	"strconv"
	"suffgo-backend-t/internal/user/entities"
	"suffgo-backend-t/internal/user/models"
	"suffgo-backend-t/internal/user/repositories"
)

type userUsecaseImpl struct {
	userRepository repositories.UserRepository
}

// en caso de no implementar la interfaz tira error
func NewUserUsecaseImpl(userRepository repositories.UserRepository) UserUsecase {
	return &userUsecaseImpl{
		userRepository: userRepository,
	}
}

func (u *userUsecaseImpl) UserDataRegister(in *models.AddUserData) error {
	insertUserData := &entities.UserDto{
		Dni:      in.Dni,
		Mail:     in.Mail,
		Password: in.Password,
		Username: in.Username,
	}

	if err := u.userRepository.InsertUserData(insertUserData); err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) GetUserByID(idStr string) (*entities.UserSafeDto, error) {
	var userData *entities.UserSafeDto
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return userData, err
	}

	userData, nil := u.userRepository.GetUserByID(id)

	return userData, nil
}

func (u *userUsecaseImpl) DeleteUser(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	if err = u.userRepository.DeleteUser(id); err != nil {
		return err
	}

	return err
}

func (u *userUsecaseImpl) GetAll() ([]entities.UserSafeDto, error) {
	users, err := u.userRepository.FetchAll()
	if err != nil {
		return nil, err
	}

	var usersDto []entities.UserSafeDto

	for _, user := range users {
		userDto := entities.UserSafeDto{
			ID:       user.ID,
			Dni:      user.Dni,
			Mail:     user.Mail,
			Username: user.Username,
		}
		usersDto = append(usersDto, userDto)
	}

	return usersDto, nil
}
