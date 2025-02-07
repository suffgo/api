package usecases

import (
	"errors"
	"suffgo/internal/settingsRoom/domain"
	e "suffgo/internal/settingsRoom/domain/errors"
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

	//tengo que fijarme si es una creacion o una actualizacion
	existingSettingRoom, err := s.repository.GetByRoom(settingRoom.RoomID())

	if errors.Is(err, e.SettingRoomNotFoundError) {
		err = s.repository.Save(settingRoom)
		if err != nil {
			return err
		}
	} else if existingSettingRoom != nil {
		return e.ErrAlreadyExists
	} else if err != nil {
		return err
	}

	return nil
}
