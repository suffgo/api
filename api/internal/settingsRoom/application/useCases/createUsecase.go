package usecases

import (
	"errors"
	"suffgo/internal/settingsRoom/domain"
	rooms "suffgo/internal/rooms/domain"
	e "suffgo/internal/settingsRoom/domain/errors"
)

type (
	CreateUsecase struct {
		repository domain.SettingRoomRepository
		roomRepository rooms.RoomRepository
	}
)

func NewCreateUsecase(repository domain.SettingRoomRepository, roomRepo rooms.RoomRepository) *CreateUsecase {
	return &CreateUsecase{
		repository: repository,
		roomRepository: roomRepo,
	}
}

func (s *CreateUsecase) Execute(settingRoom domain.SettingRoom) error {

	//solo puede tener settingroom si es formal
	room, err := s.roomRepository.GetByID(settingRoom.RoomID())

	if err != nil {
		return err
	}

	if !room.IsFormal().IsFormal {
		return errors.ErrUnsupported
	}

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
