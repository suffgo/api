package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	Vote struct {
		id       *sv.ID
		userID   *sv.ID
		optionID *sv.ID
	}

	VoteDTO struct {
		ID       uint `json:"id"`
		UserID   uint `json:"user_id"`
		OptionID uint `json:"option_id"`
	}

	VoteCreateRequest struct {
		UserID   uint `json:"user_id"`
		OptionID uint `json:"option_id"`
	}
)

func NeweVote(
	id *sv.ID,
	userID *sv.ID,
	optionID *sv.ID,
) *Vote {
	return &Vote{
		id:       id,
		userID:   userID,
		optionID: optionID,
	}
}

func (v *Vote) ID() sv.ID {
	return *v.id
}

func (v *Vote) UserID() sv.ID {
	return *v.userID
}

func (v *Vote) OptionID() sv.ID {
	return *v.optionID
}
