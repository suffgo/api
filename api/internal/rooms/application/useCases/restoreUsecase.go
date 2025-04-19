package usecases

import (
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type RestoreUsecase struct {
	roomRestoreRepository domain.RoomRepository
}

func NewRestoreUsecase(repository domain.RoomRepository) *RestoreUsecase {
	return &RestoreUsecase{
		roomRestoreRepository: repository,
	}
}

func (s *RestoreUsecase) Execute(id sv.ID) error {
	return s.roomRestoreRepository.Restore(id)
}