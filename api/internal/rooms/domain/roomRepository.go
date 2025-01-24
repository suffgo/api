package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
)

type RoomRepository interface {
	GetByID(id sv.ID) (*Room, error)
	GetAll() ([]Room, error)
	Delete(id sv.ID) error
	Save(room Room) (*Room, error)
	GetByAdminID(adminID sv.ID) ([]Room, error)
	Restore(id sv.ID) error
	SaveInviteCode(inviteCode string, roomID uint) error
	GetInviteCode(roomID uint) (string, error)
	GetRoomByCode(inviteCode string) (uint, error)
	AddToWhitelist(roomID sv.ID, userID sv.ID) error
	UserInWhitelist(roomID sv.ID, userID sv.ID) (bool, error)
}