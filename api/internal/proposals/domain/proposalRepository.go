package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
)

type ProposalRepository interface {
	GetById(id sv.ID) (*Proposal, error)
	GetAll() ([]Proposal, error)
	Delete(id sv.ID) error
	Save(proposal Proposal) (*Proposal, error)
	Update(proposal *Proposal) (*Proposal, error)
	GetByRoom(roomId sv.ID) ([]Proposal, error)
	GetResultsByRoom(roomId sv.ID) ([]ProposalResults, error)
}
