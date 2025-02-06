package usecases

import (
	"suffgo/internal/settingsRoom/domain"
	e "suffgo/internal/settingsRoom/domain/errors"
)

type UpdateSettingRoomUsecase struct {
	repository domain.SettingRoomRepository
}

func NewUpdateSettingRoomUsecase(repository domain.SettingRoomRepository) *UpdateSettingRoomUsecase {
	return &UpdateSettingRoomUsecase{
		repository: repository,
	}
}

func (u *UpdateSettingRoomUsecase) Execute(settingRoom *domain.SettingRoom) (*domain.SettingRoom, error) {
	existingSettings, err := u.repository.GetByID(settingRoom.ID())

	if err != nil {
		return nil, err
	}

	if existingSettings == nil {
		return nil, e.SettingRoomNotFoundError
	}

	updateSetting, err := u.repository.Update(settingRoom)

	if err != nil {
		return nil, err
	}

	return updateSetting, nil
}
