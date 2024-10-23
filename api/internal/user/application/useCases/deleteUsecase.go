package usecases

import "suffgo/internal/user/domain"

type DeleteUsecase struct {
	userDeleteRepository domain.UserRepository
}

func NewDeleteUsecase(repository domain.UserRepository) *DeleteUsecase {
	return &DeleteUsecase{
		userDeleteRepository: repository,
	}
}

func (s *DeleteUsecase) Execute() error {
	return nil
}