package infrastructure

import (
	"errors"
	"suffgo/cmd/database"
	d "suffgo/internal/proposals/domain"
	pe "suffgo/internal/proposals/domain/errors"
	"suffgo/internal/proposals/infrastructure/mappers"
	m "suffgo/internal/proposals/infrastructure/models"
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

func (s *ProposalXormRepository) Save(proposal d.Proposal) (*d.Proposal, error) {

	proposalModel := &m.Proposal{
		Archive:     &proposal.Archive().Archive,
		Title:       proposal.Title().Title,
		Description: &proposal.Description().Description,
		RoomID:      proposal.RoomID().Id,
	}

	_, err := s.db.GetDb().Insert(proposalModel)
	if err != nil {
		return nil, err
	}

	propMod, err := mappers.ModelToDomain(proposalModel)

	if err != nil {
		return nil, err
	}

	return propMod, nil

}

func (s *ProposalXormRepository) GetAll() ([]d.Proposal, error) {
	var proposals []m.Proposal

	err := s.db.GetDb().Where("deleted_at IS NULL").Find(&proposals)
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
		return nil, pe.ErrPropNotFound
	}

	proposalEnt, err := mappers.ModelToDomain(proposalModel)

	if err != nil {
		return nil, se.ErrDataMap
	}

	return proposalEnt, nil
}

func (s *ProposalXormRepository) Delete(id sv.ID) error {
	affected, err := s.db.GetDb().ID(id.Id).Delete(&m.Proposal{})
	if err != nil {
		return err
	}

	if affected == 0 {
		return pe.ErrPropNotFound
	}

	return nil
}

func (s *ProposalXormRepository) Restore(proposalID sv.ID) error {
	primitiveID := proposalID.Value()

	proposal := &m.Proposal{DeletedAt: nil}

	affected, err := s.db.GetDb().Unscoped().ID(primitiveID).Cols("deleted_at").Update(proposal)
	if err != nil {
		return err
	}
	if affected == 0 {
		return pe.ErrPropNotFound
	}
	return err
}

func (s *ProposalXormRepository) Update(proposal *d.Proposal) (*d.Proposal, error) {
	proposalID := proposal.ID().Id

	var existingProposal m.Proposal

	found, err := s.db.GetDb().ID(proposalID).Get(&existingProposal)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, pe.ErrPropNotFound
	}

	updateProposal := mappers.DomainToModel(proposal)

	affected, err := s.db.GetDb().ID(proposalID).Update(updateProposal)

	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, errors.New("no rows were updated")
	}

	updatedProposal, err := mappers.ModelToDomain(updateProposal)

	if err != nil {
		return nil, err
	}

	return updatedProposal, nil
}

func (s *ProposalXormRepository) GetByRoom(roomId sv.ID) ([]d.Proposal, error) {

	var proposals []m.Proposal

	err := s.db.GetDb().Where("deleted_at IS NULL and room_id = ?", roomId.Id).Find(&proposals)
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

func (s *ProposalXormRepository) GetResultsByRoom(roomId sv.ID) ([]d.ProposalResults, error) {
	var rawResults []m.SqlResult
	err := s.db.GetDb().SQL(`
		SELECT 
			p.id AS proposal_id,
			p.title AS proposal_title,
			p.description AS proposal_description,
			o.id AS option_id,
			o.value AS option_value,
			v.id AS vote_id,
			u.id AS user_id,
			u.username AS username,
			u.image AS user_image
		FROM proposal p
		LEFT JOIN option o ON o.proposal_id = p.id
		LEFT JOIN vote v ON v.option_id = o.id
		LEFT JOIN users u ON u.id = v.user_id
		WHERE p.room_id = ?
		ORDER BY p.id, o.id, v.id
	`, &roomId.Id).Find(&rawResults)

	if err != nil {
		return nil, err
	}

	// Procesar resultados
	var results []d.ProposalResults
	var currentProposal *d.ProposalResults
	var currentOption *d.OptionResults

	for _, row := range rawResults {
		if currentProposal == nil || currentProposal.ProposalId != row.ProposalId {
			if currentProposal != nil {
				if currentOption != nil {
					currentProposal.Options = append(currentProposal.Options, *currentOption)
				}
				results = append(results, *currentProposal)
			}

			currentProposal = &d.ProposalResults{
				ProposalId:          row.ProposalId,
				ProposalTitle:       row.ProposalTitle,
				ProposalDescription: row.ProposalDescription,
				Options:             []d.OptionResults{},
			}
			currentOption = nil
		}

		if row.OptionId != 0 && (currentOption == nil || currentOption.OptionId != row.OptionId) {
			if currentOption != nil {
				currentProposal.Options = append(currentProposal.Options, *currentOption)
			}
			currentOption = &d.OptionResults{
				OptionId:    row.OptionId,
				OptionValue: row.OptionValue,
				Votes:       []d.VotesResults{},
			}
		}

		if row.VoteId != 0 && currentOption != nil {
			currentOption.Votes = append(currentOption.Votes, d.VotesResults{
				VoteId:    row.VoteId,
				UserId:    row.UserId,
				Username:  row.Username,
				UserImage: row.UserImage,
			})
		}
	}

	if currentProposal != nil {
		if currentOption != nil {
			currentProposal.Options = append(currentProposal.Options, *currentOption)
		}
		results = append(results, *currentProposal)
	}

	return results, nil
}
