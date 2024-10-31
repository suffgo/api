package domain

import (
	v "suffgo/internal/option/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	Option struct {
		id         *sv.ID
		value      v.Value
		proposalID *sv.ID
	}

	OptionDTO struct {
		ID         uint   `json:"id"`
		Value      string `json:"value"`
		ProposalID uint   `json:"proposal_id"`
	}

	OptionCreateRequest struct {
		Value      string `json:"value"`
		ProposalID uint   `json:"proposal_id"`
	}
)

func NewOption(
	id *sv.ID,
	value v.Value,
	proposalID *sv.ID,
) *Option {
	return &Option{
		id:         id,
		value:      value,
		proposalID: proposalID,
	}
}

func (o *Option) ID() sv.ID {
	return *o.id
}

func (o *Option) Value() v.Value {
	return o.value
}

func (o *Option) ProposalID() sv.ID {
	return *o.proposalID
}
