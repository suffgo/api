package infrastructure

import (
	"suffgo/cmd/database"
	d "suffgo/internal/proposal/domain"
	m "suffgo/internal/proposal/infrastructure/models"
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
	archiveValue := proposal.Archive().Archive // Esto es un string
	descriptionValue := proposal.Description().Description

	proposalModel := &m.Proposal{
		Archive:     &archiveValue,
		Title:       proposal.Title().Title,
		Description: &descriptionValue,
	}

	_, err := s.db.GetDb().Insert(proposalModel)
	if err != nil {
		return err
	}

	return nil

}
