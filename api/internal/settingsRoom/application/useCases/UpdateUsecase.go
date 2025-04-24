package usecases

import (
	"errors"
	rd "suffgo/internal/rooms/domain"
	"suffgo/internal/settingsRoom/domain"
	e "suffgo/internal/settingsRoom/domain/errors"
)

type UpdateSettingRoomUsecase struct {
	settingRepository domain.SettingRoomRepository
	roomRepository    rd.RoomRepository
}

func NewUpdateSettingRoomUsecase(settingRepo domain.SettingRoomRepository, roomRepo rd.RoomRepository) *UpdateSettingRoomUsecase {
	return &UpdateSettingRoomUsecase{
		settingRepository: settingRepo,
		roomRepository:    roomRepo,
	}
}

func (u *UpdateSettingRoomUsecase) Execute(settingRoom *domain.SettingRoom) (*domain.SettingRoom, error) {
	existingSettings, err := u.settingRepository.GetByID(settingRoom.ID())
	if err != nil {
		return nil, err
	}
	room, err := u.roomRepository.GetByID(settingRoom.RoomID())
	if err != nil {
		return nil, err
	}

	if room.State().CurrentState == "online" {
		return nil, errors.New("La sala esta activa, no se puede modificar")
	}

	if !room.IsFormal().IsFormal {
		return nil, errors.New("No se puede configurar salas informales")
	}

	if existingSettings == nil {
		return nil, e.SettingRoomNotFoundError
	}

	updateSetting, err := u.settingRepository.Update(settingRoom)

	if err != nil {
		return nil, err
	}

	return updateSetting, nil
}
