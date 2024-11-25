package usecases

import (
	"suffgo/internal/proposals/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	DeleteUseCase struct {
		deleteRepository domain.ProposalRepository
	}
)

func NewDeleteUseCase(repository domain.ProposalRepository) *DeleteUseCase {
	return &DeleteUseCase{
		deleteRepository: repository,
	}
}

func (s *DeleteUseCase) Execute(id sv.ID) error {
	err := s.deleteRepository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
