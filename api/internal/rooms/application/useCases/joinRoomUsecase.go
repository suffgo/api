package usecases

import (
	"errors"
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
	sr "suffgo/internal/rooms/domain/valueObjects"
	rerr "suffgo/internal/rooms/domain/errors"
)

type JoinRoomUsecase struct {
	joinRoomUsecaseRepository domain.RoomRepository
	
}

func NewJoinRoomUsecase(repository domain.RoomRepository) *JoinRoomUsecase {
	return &JoinRoomUsecase{
		joinRoomUsecaseRepository: repository,
	}
}

//Metodo importante, gestiona la union del usuario a la sala
func (s *JoinRoomUsecase) Execute(roomCode string, userID sv.ID) (*domain.Room, error) {
	//Obtener sala a traves de codigo
	roomID, err := s.joinRoomUsecaseRepository.GetRoomByCode(roomCode)
	
	if err != nil {
		return nil, err
	}

	rID := sv.ID{Id: roomID}
	//obtener datos de sala
	room, err := s.joinRoomUsecaseRepository.GetByID(rID)

	if err != nil {
		return nil, errors.New("error al obtener la sala")
	}

	if room == nil {
		return nil, errors.New("error al obtener la sala")
	}
	//Si la sala es formal verificar si el usuario puede unirse a la misma (tiene permiso, cantidad maxima, etc )
	if room.IsFormal().IsFormal {
		//check whitelist en user_room. Aca estoy asumiendo que todas las salas formales usan whitelist
		can, err := s.joinRoomUsecaseRepository.UserInWhitelist(room.ID(), userID)

		if err != nil {
			return nil, err
		}

		if !can {
			return nil, rerr.ErrNotWhitelist
		}
	}

	code := sr.InviteCode{Code: roomCode}
	room.SetInviteCode(code)
	return room, nil
}