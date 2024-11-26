package usecases

import (
	"errors"
	"suffgo/internal/settingsRoom/domain"
)

type (
	CreateUsecase struct {
		repository domain.SettingRoomRepository
	}
)

func NewCreateUsecase(repository domain.SettingRoomRepository) *CreateUsecase {
	return &CreateUsecase{
		repository: repository,
	}
}

func (s *CreateUsecase) Execute(settingRoom domain.SettingRoom) error {
	existingSettingRoom, err := s.repository.GetByID(settingRoom.ID())

	if err != nil {
		return err
	}

	if existingSettingRoom != nil {
		return errors.New("opcion ya existe")
	}

	err = s.repository.Save(settingRoom)
	if err != nil {
		return err
	}
	return nil
}
