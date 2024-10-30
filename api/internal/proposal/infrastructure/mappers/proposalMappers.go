package mappers

import (
	"suffgo/internal/proposal/domain"
	v "suffgo/internal/proposal/domain/valueObjects"
	m "suffgo/internal/proposal/infrastructure/models"
	sv "suffgo/internal/shared/domain/valueObjects"
)

func DomainToModel(proposal *domain.Proposal) *m.Proposal {
	return &m.Proposal{
		ID:          proposal.ID().Id, // Convierte ID a uint
		Archive:     &proposal.Archive().Archive,
		Title:       proposal.Title().Title,
		Description: &proposal.Description().Description,
	}
}

func ModelToDomain(proposalModel *m.Proposal) (*domain.Proposal, error) {
	id, err := sv.NewID(proposalModel.ID)
	if err != nil {
		return nil, err
	}

	archive, err := v.NewArchive(*proposalModel.Archive)
	if err != nil {
		return nil, err
	}

	title, err := v.NewTitle(proposalModel.Title)
	if err != nil {
		return nil, err
	}

	description, err := v.NewDescription(*proposalModel.Description)
	if err != nil {
		return nil, err
	}

	return domain.NewProposal(id, archive, *title, description), nil
}
