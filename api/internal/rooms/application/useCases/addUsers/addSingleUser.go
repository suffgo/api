package addusers

import (
	"errors"
	"suffgo/internal/rooms/domain"
	roomErrors "suffgo/internal/rooms/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
	userDomain "suffgo/internal/users/domain"
	userErr "suffgo/internal/users/domain/errors"
	userDomainV "suffgo/internal/users/domain/valueObjects"
)

type AddSingleUserUsecase struct {
	repository     domain.RoomRepository
	userRepository userDomain.UserRepository
}

func NewAddSingleUserUsecase(repository domain.RoomRepository, userRepository userDomain.UserRepository) *AddSingleUserUsecase {
	return &AddSingleUserUsecase{
		repository:     repository,
		userRepository: userRepository,
	}
}

func (s *AddSingleUserUsecase) Execute(userData string, roomID, adminID sv.ID) error {

	//chequear que el administrador de la sala sea el que esta intentando agregar usuarios
	room, err := s.repository.GetByID(roomID)

	if err != nil {
		return err
	}

	if room == nil {
		return roomErrors.ErrRoomNotFound
	}

	if room.AdminID().Id != adminID.Id {
		return roomErrors.ErrUserNotAdmin
	}

	if room.State().CurrentState == "online" {
		return errors.New("La sala esta activa, no se puede agregar nuevos usuarios")
	}

	user, err := s.lookForUser(userData)

	if err != nil {
		return err
	}

	already, err := s.repository.UserInWhitelist(roomID, user.ID())

	if already {
		return roomErrors.ErrAlreadyInWhitelist
	}

	if user != nil {

		err = s.repository.AddToWhitelist(roomID, user.ID())

		if err != nil {
			return nil
		}

	}

	return nil
}

func (s *AddSingleUserUsecase) lookForUser(userData string) (*userDomain.User, error) {
	//Tengo que ver que tipo de dato es
	var user *userDomain.User
	//Si es mail
	mail, err := userDomainV.NewEmail(userData)
	if err == nil {
		//Obtengo user por mail
		user, err = s.userRepository.GetByEmail(*mail)

		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}

	//Si es nombre de usuario
	username, err := userDomainV.NewUserName(userData)
	if err == nil {
		//Obtengo user por nombre de usuario
		user, err = s.userRepository.GetByUsername(*username)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}

	dni, err := userDomainV.NewDni(userData)
	if err == nil {
		//Obtengo user por dni
		user, err = s.userRepository.GetByDni(*dni)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return user, nil
		}
	}

	return nil, userErr.ErrUserNotFound
}
