package usecases

import (
	"suffgo/internal/votes/domain"
)

type (
	CreateUsecase struct {
		repository domain.VoteRepository
	}
)

func NewCreateUsecase(repository domain.VoteRepository) *CreateUsecase {
	return &CreateUsecase{
		repository: repository,
	}
}

func (s *CreateUsecase) Execute(vote domain.Vote) error {
	return nil
}
