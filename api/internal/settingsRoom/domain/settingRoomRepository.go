package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
)

type SettingRoomRepository interface {
	GetByID(id sv.ID) (*SettingRoom, error)
	GetAll() ([]SettingRoom, error)
	Delete(id sv.ID) error
	Save(settingRoom SettingRoom) error
	//Update(settingRoom SettingRoom) error
}
