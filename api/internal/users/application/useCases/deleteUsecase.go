package usecases

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	"suffgo/internal/users/domain"
)

type DeleteUsecase struct {
	userDeleteRepository domain.UserRepository
}

func NewDeleteUsecase(repository domain.UserRepository) *DeleteUsecase {
	return &DeleteUsecase{
		userDeleteRepository: repository,
	}
}

func (s *DeleteUsecase) Execute(id sv.ID) error {

	err := s.userDeleteRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
