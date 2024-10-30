package infrastructure

import (
	"suffgo/cmd/database"
	d "suffgo/internal/proposal/domain"
	ue "suffgo/internal/proposal/domain/errors"
	"suffgo/internal/proposal/infrastructure/mappers"
	m "suffgo/internal/proposal/infrastructure/models"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type ProposalXormRepository struct {
	db database.Database
}

func NewProposalXormRepository(db database.Database) *ProposalXormRepository {
	return &ProposalXormRepository{
		db: db,
	}
}

func (s *ProposalXormRepository) Save(proposal d.Proposal) error {

	proposalModel := &m.Proposal{
		Archive:     &proposal.Archive().Archive,
		Title:       proposal.Title().Title,
		Description: &proposal.Description().Description,
	}

	_, err := s.db.GetDb().Insert(proposalModel)
	if err != nil {
		return err
	}

	return nil

}

func (s *ProposalXormRepository) GetAll() ([]d.Proposal, error) {
	var proposals []m.Proposal

	err := s.db.GetDb().Find(&proposals)
	if err != nil {
		return nil, err
	}

	var proposalsDomain []d.Proposal
	for _, proposal := range proposals {
		proposalDomain, err := mappers.ModelToDomain(&proposal)

		if err != nil {
			return nil, err
		}

		proposalsDomain = append(proposalsDomain, *proposalDomain)
	}
	return proposalsDomain, nil
}

func (s *ProposalXormRepository) GetById(id sv.ID) (*d.Proposal, error) {
	proposalModel := new(m.Proposal)
	has, err := s.db.GetDb().ID(id.Id).Get(proposalModel)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ue.ProposalNotFoundError
	}

	proposalEnt, err := mappers.ModelToDomain(proposalModel)

	if err != nil {
		return nil, se.DataMappingError
	}

	return proposalEnt, nil
}

func (s *ProposalXormRepository) Delete(id sv.ID) error {
	affected, err := s.db.GetDb().ID(id.Id).Delete(&m.Proposal{})
	if err != nil {
		return err
	}

	if affected == 0 {
		return ue.ProposalNotFoundError
	}

	return nil
}
