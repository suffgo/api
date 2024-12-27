package usecases

import (
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type RestoreUsecase struct {
	userRestoreRepository domain.RoomRepository
}

func NewRestoreUsecase(repository domain.RoomRepository) *RestoreUsecase {
	return &RestoreUsecase{
		userRestoreRepository: repository,
	}
}

func (s *RestoreUsecase) Execute(id sv.ID) error {
	return s.userRestoreRepository.Restore(id)
}
