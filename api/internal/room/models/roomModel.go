package models

import (
	ur "suffgo/internal/user_room/entities"
)

type AddRoomData struct {
	LinkInvite   *string
	IsFormal     bool
	Name         string
	AdminID      uint
	AllowedUsers []ur.UserRoom
}
