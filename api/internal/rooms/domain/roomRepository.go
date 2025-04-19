package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
)

type RoomRepository interface {
	GetByID(id sv.ID) (*Room, error)
	GetAll() ([]Room, error)
	Delete(roomID sv.ID) error
	Save(room Room) (*Room, error)
	GetByAdminID(adminID sv.ID) ([]Room, error)
	Restore(id sv.ID) error
	GetRoomByCode(inviteCode string) (*Room, error)
	AddToWhitelist(roomID sv.ID, userID sv.ID) error
	UserInWhitelist(roomID sv.ID, userID sv.ID) (bool, error)
	Update(room *Room) (*Room, error)
	RemoveFromWhitelist(roomId sv.ID, userId sv.ID) error
	RestartRoom(roomId sv.ID) error
	HistoryRooms(userId sv.ID) ([]Room, error)
}
