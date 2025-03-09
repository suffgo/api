package usecases

import (
	opterr "suffgo/internal/options/domain/errors"
	"suffgo/internal/options/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	GetByPropUsecase struct {
		repository domain.OptionRepository
	}
)

func NewGetByPropUsecase(repository domain.OptionRepository) *GetByPropUsecase {
	return &GetByPropUsecase{
		repository: repository,
	}
}

func (s *GetByPropUsecase) Execute(proposalId  *sv.ID) ([]domain.Option, error) {

	options, err := s.repository.GetByProposal(*proposalId)

	if err != nil {
		return nil, err
	}

	if len(options) == 0 {
		return nil, opterr.ErrOptNotFound
	}


	return options, nil
}
