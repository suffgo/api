package usecases

import "suffgo/internal/users/domain"

type CreateSessionUsecase struct {
	repository domain.UserRepository
}

func NewCreateSessionUsecase(repo domain.UserRepository) *CreateSessionUsecase {
	return &CreateSessionUsecase{
		repository: repo,
	}
}

func (s *CreateSessionUsecase) Execute() error {
	return nil
}
