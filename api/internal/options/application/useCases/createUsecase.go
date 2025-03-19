package usecases

import (
	"suffgo/internal/options/domain"
	opterr "suffgo/internal/options/domain/errors"
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

//Por ahora se maneja la duplicacion de opciones en el frontend
func (s *CreateUsecase) Execute(option domain.Option) error {
	err := s.repository.Save(option)
	if err != nil {
		return opterr.ErrOptRepeated
	}
	return nil
}
