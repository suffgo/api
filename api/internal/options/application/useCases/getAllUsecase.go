package usecases

import "suffgo/internal/options/domain"

type GetAllUsecase struct {
	optionGetAllRepository domain.OptionRepository
}

func NewGetAllRepository(repository domain.OptionRepository) *GetAllUsecase {
	return &GetAllUsecase{
		optionGetAllRepository: repository,
	}
}

func (s *GetAllUsecase) Execute() ([]domain.Option, error) {
	options, err := s.optionGetAllRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return options, nil
}
