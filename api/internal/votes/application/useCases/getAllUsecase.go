package usecases

import (
	"suffgo/internal/votes/domain"
)

type GetAllUsecase struct {
	voteGetAllRepository domain.VoteRepository
}

func NewGetAllRepository(repository domain.VoteRepository) *GetAllUsecase {
	return &GetAllUsecase{
		voteGetAllRepository: repository,
	}
}

func (s *GetAllUsecase) Execute() ([]domain.Vote, error) {
	votes, err := s.voteGetAllRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return votes, nil
}
