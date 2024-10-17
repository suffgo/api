package entities

import (
	ur "suffgo/internal/user_room/entities"

	"gorm.io/gorm"
)

type (
	Room struct {
		gorm.Model
		LinkInvite   *string
		IsFormal     bool
		Name         string
		AdminID      uint
		AllowedUsers []ur.UserRoom `gorm:"foreignKey:RoomID"`
	}

	RoomDto struct {
		ID           uint          `json:"id"`
		LinkInvite   *string       `json:"linkInvite"`
		IsFormal     bool          `json:"isFormal"`
		Name         string        `json:"name"`
		AdminID      uint          `json:"admin"`
		AllowedUsers []ur.UserRoom `json:"allowedUsers"`
	}
)
