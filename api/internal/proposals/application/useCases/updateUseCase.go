package usecases

import (
	d "suffgo/internal/proposals/domain"
	e "suffgo/internal/proposals/domain/errors"
)

type UpdateUsecase struct {
	repository d.ProposalRepository
}

func NewUpdateProposalUsecase(repository d.ProposalRepository) *UpdateUsecase {
	return &UpdateUsecase{
		repository: repository,
	}
}

func (u *UpdateUsecase) Execute(proposal *d.Proposal) (*d.Proposal, error) {
	existingProposal, err := u.repository.GetById(proposal.ID())

	if err != nil {
		return nil, err
	}

	if existingProposal == nil {
		return nil, e.ErrPropNotFound
	}

	updateProposal, err := u.repository.Update(proposal)

	if err != nil {
		return nil, err
	}

	return updateProposal, nil
}
