package usecases

import (
	"errors"
	"suffgo/internal/option/domain"
)

type (
	CreateUsecase struct {
		repository domain.OptionRepository
	}
)

func NewCreateUsecase(repository domain.OptionRepository) *CreateUsecase {
	return &CreateUsecase{
		repository: repository,
	}
}

func (s *CreateUsecase) Execute(option domain.Option) error {
	existingOption, err := s.repository.GetByValue(option.Value())

	if err != nil {
		return err
	}
	if existingOption != nil {
		return errors.New("opcion ya existe")
	}

	err = s.repository.Save(option)
	if err != nil {
		return err
	}
	return nil
}
