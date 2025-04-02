package usecases

import (
	"errors"
	roomDom "suffgo/internal/rooms/domain"
	setrDom "suffgo/internal/settingsRoom/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
	"suffgo/internal/users/domain"
)

type GetUsersByRoom struct {
	usrRep domain.UserRepository
	roomRep roomDom.RoomRepository
	setrRep setrDom.SettingRoomRepository
}

func NewGetUsersByRoom(usrRep domain.UserRepository, roomRep roomDom.RoomRepository, setrRep setrDom.SettingRoomRepository) *GetUsersByRoom {
	return &GetUsersByRoom{
		usrRep: usrRep,
		roomRep: roomRep,
		setrRep: setrRep,
	}
}

func (s *GetUsersByRoom) Execute(roomId sv.ID) ([]domain.User, error) {
	
	//obtengo la sala
	room, err := s.roomRep.GetByID(roomId)

	if err != nil {
		return nil, err
	}

	//necesito saber si es privada o publica asi que obtengo su settingRoom
	roomSetr, err := s.setrRep.GetByRoom(room.ID())

	if err != nil {
		return nil, err
	}

	if *roomSetr.Privacy().Privacy {
		users, err := s.usrRep.GetByRoom(room.ID())

		if err != nil {
			return nil, err
		}

		return users, nil
	} else {
		return nil, errors.New("room is public")
	}	
}
