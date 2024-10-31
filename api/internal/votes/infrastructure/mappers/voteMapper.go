package mappers

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	"suffgo/internal/votes/domain"
	m "suffgo/internal/votes/infrastructure/models"
)

func DomainToModel(vote *domain.Vote) *m.Vote {
	return &m.Vote{
		ID:       vote.ID().Id,
		UserID:   vote.UserID().Id,
		OptionID: vote.OptionID().Id,
	}
}

func ModelToDomain(voteModel *m.Vote) (*domain.Vote, error) {
	id, err := sv.NewID(voteModel.ID)
	if err != nil {
		return nil, err
	}
	userID, err := sv.NewID(voteModel.UserID)
	if err != nil {
		return nil, err
	}
	proposalID, err := sv.NewID(voteModel.OptionID)
	if err != nil {
		return nil, err
	}
	return domain.NeweVote(id, userID, proposalID), nil
}
