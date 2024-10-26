package usecases

import (
	"suffgo/internal/option/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type DeleteUsecase struct {
	optionDeleteRepository domain.OptionRepository
}

func NewDeleteUsecase(repository domain.OptionRepository) *DeleteUsecase {
	return &DeleteUsecase{
		optionDeleteRepository: repository,
	}
}

func (s *DeleteUsecase) Execute(id sv.ID) error {
	err := s.optionDeleteRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
