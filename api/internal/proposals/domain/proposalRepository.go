package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
)

type ProposalRepository interface {
	GetById(id sv.ID) (*Proposal, error)
	GetAll() ([]Proposal, error)
	Delete(id sv.ID) error
	//Edit()
	Save(proposal Proposal) (*Proposal, error)
}
