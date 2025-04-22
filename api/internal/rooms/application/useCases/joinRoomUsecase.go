package usecases

import (
	"errors"
	"suffgo/internal/rooms/domain"
	rerr "suffgo/internal/rooms/domain/errors"
	srdom "suffgo/internal/settingsRoom/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type JoinRoomUsecase struct {
	roomRepo        domain.RoomRepository
	setrRepo srdom.SettingRoomRepository
}

func NewJoinRoomUsecase(repository domain.RoomRepository, srRepo srdom.SettingRoomRepository) *JoinRoomUsecase {
	return &JoinRoomUsecase{
		roomRepo: repository,
		setrRepo: srRepo,
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

	setroom, err := s.setrRepo.GetByRoom(room.ID())

	if err != nil {
		return nil, err
	}


	if *setroom.Privacy().Privacy {
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
