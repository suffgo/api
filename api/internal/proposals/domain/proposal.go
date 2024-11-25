package domain

import (
	v "suffgo/internal/proposals/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	Proposal struct {
		id          *sv.ID
		archive     *v.Archive
		title       v.Title
		description *v.Description
		roomID      sv.ID
	}

	ProposalDTO struct {
		ID          uint    `json:"id"`
		Archive     *string `json:"archive"`
		Title       string  `json:"title"`
		Description *string `json:"description"`
		RoomID      uint    `json:"room_id"`
	}

	ProposalCreateRequest struct {
		Archive       *string `json:"archive"`
		Title         string  `json:"title"`
		Description   *string `json:"description"`
		RoomID        uint    `json:"room_id"`
		UserCreatorID uint    `json:"user_creator_id"`
	}
)

func NewProposal(
	id *sv.ID,
	archive *v.Archive,
	title v.Title,
	description *v.Description,
	roomID *sv.ID,
) *Proposal {
	return &Proposal{
		id:          id,
		archive:     archive,
		title:       title,
		description: description,
		roomID:      *roomID,
	}
}

func (p *Proposal) ID() sv.ID {
	return *p.id
}

func (p *Proposal) Archive() *v.Archive {
	return p.archive
}

func (p *Proposal) Title() v.Title {
	return p.title
}

func (p *Proposal) Description() *v.Description {
	return p.description
}

func (p *Proposal) RoomID() sv.ID {
	return p.roomID
}
