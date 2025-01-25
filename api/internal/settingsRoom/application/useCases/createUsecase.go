package usecases

import (
	e "suffgo/internal/settingsRoom/domain/errors"
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
	
	//tengo que ver si la existe 


	//tengo que fijarme si es una creacion o una actualizacion
	existingSettingRoom, err := s.repository.GetByRoom(settingRoom.RoomID())

	if err != nil {
		return err
	}

	if existingSettingRoom != nil {
		return e.ErrAlreadyExists
	}

	err = s.repository.Save(settingRoom)
	if err != nil {
		return err
	}
	return nil
}
