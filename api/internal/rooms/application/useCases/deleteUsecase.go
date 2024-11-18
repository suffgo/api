package usecases

import (
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

func (s *DeleteUsecase) Execute(id sv.ID) error {

	err := s.roomDeleteRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
