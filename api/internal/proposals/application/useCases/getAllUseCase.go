package usecases

import (
	"suffgo/internal/proposals/domain"
)

type (
	GetAllUsecase struct {
		getAllRepository domain.ProposalRepository
	}
)

func NewGetAllUseCase(repository domain.ProposalRepository) *GetAllUsecase {

	return &GetAllUsecase{
		getAllRepository: repository,
	}

}

func (s *GetAllUsecase) Execute() ([]domain.Proposal, error) {

	proposal, err := s.getAllRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return proposal, nil
}
