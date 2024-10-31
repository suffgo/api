package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
)

type VoteRepository interface {
	GetByID(id sv.ID) (*Vote, error)
	GetAll() ([]Vote, error)
	Delete(id sv.ID) error
	Save(vote Vote) error
}
