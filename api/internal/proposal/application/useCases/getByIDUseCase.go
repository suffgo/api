package usecases

import (
	"suffgo/internal/proposal/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	GetByIDUsecase struct {
		getByIdRepository domain.ProposalRepository
	}
)

func NewGetByIDUseCase(repository domain.ProposalRepository) *GetByIDUsecase {
	return &GetByIDUsecase{
		getByIdRepository: repository,
	}

}

func (s *GetByIDUsecase) Execute(id sv.ID) (*domain.Proposal, error) {
	proposal, err := s.getByIdRepository.GetById(id)

	if err != nil {
		return nil, err
	}

	return proposal, nil
}
