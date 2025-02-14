package usecases

import (
	"errors"
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type DeleteUsecase struct {
	roomDeleteRepository domain.RoomRepository
}

func NewDeleteUsecase(repository domain.RoomRepository) *DeleteUsecase {
	return &DeleteUsecase{
		roomDeleteRepository: repository,
	}
}

func (s *DeleteUsecase) Execute(roomID sv.ID, userID sv.ID) error {

	room, err := s.roomDeleteRepository.GetByID(roomID)
	if err != nil {
		return err
	}

	if room.AdminID() != userID {
		return errors.New("unauthorized")
	}

	err = s.roomDeleteRepository.Delete(roomID)

	if err != nil {
		return err
	}

	return nil
}
