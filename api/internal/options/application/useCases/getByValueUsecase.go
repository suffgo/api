package usecases

import (
	"suffgo/internal/options/domain"
	v "suffgo/internal/options/domain/valueObjects"
)

type GetByValueUsecase struct {
	optionGetByValueRepository domain.OptionRepository
}

func NewGetByValueUsecase(repository domain.OptionRepository) *GetByValueUsecase {
	return &GetByValueUsecase{
		optionGetByValueRepository: repository,
	}
}

func (s *GetByValueUsecase) Execute(value v.Value) (*domain.Option, error) {
	option, err := s.optionGetByValueRepository.GetByValue(value)

	if err != nil {
		return nil, err
	}
	return option, nil
}
