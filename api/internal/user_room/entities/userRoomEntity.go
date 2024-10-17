package entities

import "gorm.io/gorm"

type (
	UserRoom struct {
		gorm.Model
		UserID uint
		RoomID uint
	}
)
