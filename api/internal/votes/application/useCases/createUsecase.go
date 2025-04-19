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

func (s *CreateUsecase) Execute(voteData domain.Vote) (*domain.Vote, error) {
	createdVote, err := s.repository.Save(voteData)
	if err != nil {
		return nil, err
	}

	return createdVote, nil
}
