package usecases

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	"suffgo/internal/votes/domain"
)

type GetByIDUsecase struct {
	voteGetByIDReposiroty domain.VoteRepository
}

func NewGetByIDUsecase(repository domain.VoteRepository) *GetByIDUsecase {
	return &GetByIDUsecase{
		voteGetByIDReposiroty: repository,
	}
}

func (s *GetByIDUsecase) Execute(id sv.ID) (*domain.Vote, error) {
	vote, err := s.voteGetByIDReposiroty.GetByID(id)

	if err != nil {
		return nil, err
	}

	return vote, nil

}
