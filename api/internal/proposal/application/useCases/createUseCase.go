package usecases

import (
	"suffgo/internal/proposal/domain"
)

type (
	CreateUsecase struct {
		repository domain.ProposalRepository
	}
)

func NewCreateUsecase(repository domain.ProposalRepository) *CreateUsecase {
	return &CreateUsecase{
		repository: repository,
	}
}

func (s *CreateUsecase) Execute(proposal domain.Proposal) error {

	err := s.repository.Save(proposal)

	if err != nil {
		return err
	}

	return nil

}
