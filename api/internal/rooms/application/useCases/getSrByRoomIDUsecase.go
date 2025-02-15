package usecases

import (
	"suffgo/internal/rooms/domain"
	sr "suffgo/internal/settingsRoom/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type GetSrByRoomUsecase struct {
	roomRepository        domain.RoomRepository
	settingRoomRepository sr.SettingRoomRepository
}

func NewGetSrByRoomUsecase(roomRepo domain.RoomRepository, srRepo sr.SettingRoomRepository) *GetSrByRoomUsecase {
	return &GetSrByRoomUsecase{
		roomRepository:        roomRepo,
		settingRoomRepository: srRepo,
	}
}

func (s *GetSrByRoomUsecase) Execute(roomID sv.ID) (*sr.SettingRoom, error) {

	settingRoom, err :=  s.settingRoomRepository.GetByRoom(roomID)

	if err != nil {
		return nil, err
	}


	return settingRoom, nil
}
