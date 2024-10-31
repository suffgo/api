package usecases

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	"suffgo/internal/votes/domain"
)

type DeleteUsecase struct {
	voteDeleteRepository domain.VoteRepository
}

func NewDeleteUsecase(repository domain.VoteRepository) *DeleteUsecase {
	return &DeleteUsecase{
		voteDeleteRepository: repository,
	}
}

func (s *DeleteUsecase) Execute(id sv.ID) error {

	err := s.voteDeleteRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
