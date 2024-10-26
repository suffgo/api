package usecases

import (
	"suffgo/internal/option/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type GetByIDUsecase struct {
	optionGetByIDRespository domain.OptionRepository
}

func NewGetByIDUsecase(repository domain.OptionRepository) *GetByIDUsecase {
	return &GetByIDUsecase{
		optionGetByIDRespository: repository,
	}
}

func (s *GetByIDUsecase) Execute(id sv.ID) (*domain.Option, error) {

	option, err := s.optionGetByIDRespository.GetByID(id)

	if err != nil {
		return nil, err
	}
	return option, nil
}
