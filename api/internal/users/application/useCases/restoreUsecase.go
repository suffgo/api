package usecases

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	"suffgo/internal/users/domain"
)

type RestoreUsecase struct {
	userRestoreRepository domain.UserRepository
}

func NewRestoreUsecase(repository domain.UserRepository) *RestoreUsecase {
	return &RestoreUsecase{
		userRestoreRepository: repository,
	}
}

func (s *RestoreUsecase) Execute(id sv.ID) error {
	return s.userRestoreRepository.Restore(id)
}
