package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
)

type RoomRepository interface {
	GetByID(id sv.ID) (*Room, error)
	GetAll() ([]Room, error)
	Delete(id sv.ID) error
	Save(room Room) error
	GetByAdmin(adminID sv.ID) ([]Room, error)
}
