package usecases

import (
	"errors"
	"suffgo/internal/rooms/domain"
	rerr "suffgo/internal/rooms/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type JoinRoomUsecase struct {
	roomRepo domain.RoomRepository
}

func NewJoinRoomUsecase(repository domain.RoomRepository) *JoinRoomUsecase {
	return &JoinRoomUsecase{
		roomRepo: repository,
	}
}

func (s *JoinRoomUsecase) Execute(roomCode string, userID sv.ID) (*domain.Room, error) {
	room, err := s.roomRepo.GetRoomByCode(roomCode)

	if err != nil {
		return nil, err
	}

	if room == nil {
		return nil, errors.New("error al obtener la sala")
	}

	if room.IsFormal().IsFormal {
		//check whitelist en user_room. Aca estoy asumiendo que todas las salas formales usan whitelist
		can, err := s.roomRepo.UserInWhitelist(room.ID(), userID)

		if err != nil {
			return nil, err
		}

		if !can {
			return nil, rerr.ErrNotWhitelist
		}
	}

	return room, nil
}
