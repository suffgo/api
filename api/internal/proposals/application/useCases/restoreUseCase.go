package usecases

import (
	"suffgo/internal/proposals/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type RestoreUsecase struct {
	proposalRestoreRepository domain.ProposalRepository
}

func NewRestoreUsecase(repository domain.ProposalRepository) *RestoreUsecase {
	return &RestoreUsecase{
		proposalRestoreRepository: repository,
	}
}

func (s *RestoreUsecase) Execute(id sv.ID) error {
	return s.proposalRestoreRepository.Restore(id)
}
