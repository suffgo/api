package domain

import (
	v "suffgo/internal/proposal/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	Proposal struct {
		id          *sv.ID
		archive     v.Archive
		title       v.Title
		description v.Description
	}

	PorposalDTO struct {
		ID          uint    `json:"id"`
		Archive     *string `json:"archive"`
		Title       string  `json:"title"`
		Description *string `json:"description"`
	}

	PorposalCreateRequest struct {
		Archive     *string `json:"archive"`
		Title       string  `json:"title"`
		Description *string `json:"description"`
	}
)

func NewProposal(
	id *sv.ID,
	archive v.Archive,
	title v.Title,
	description v.Description,

) *Proposal {
	return &Proposal{
		id:          id,
		archive:     archive,
		title:       title,
		description: description,
	}
}

func (u *Proposal) ID() sv.ID {
	return *u.id
}

func (u *Proposal) Archive() v.Archive {
	return u.archive
}

func (u *Proposal) Title() v.Title {
	return u.title
}

func (u *Proposal) Description() v.Description {
	return u.description
}
