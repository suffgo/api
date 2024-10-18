package entities

import (
	r "suffgo/internal/room/entities"
	ur "suffgo/internal/user_room/entities"

	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Dni          string        `gorm:"unique; not null"`
		Mail         string        `gorm:"unique; not null"`
		Password     string        `gorm:"not null"`
		Username     string        `gorm:"unique; not null"`
		CreatedRooms []r.Room      `gorm:"foreignKey:AdminID"`
		AllowedRooms []ur.UserRoom `gorm:"foreignKey:UserID"` // salas a las que puede ingresar
	}

	UserDto struct {
		ID           uint     `json:"id"`
		Dni          string   `json:"dni"`
		Mail         string   `json:"mail"`
		Password     string   `json:"password"`
		Username     string   `json:"username"`
		CreatedRooms []r.Room `json:"createdRooms"`
	}

	UserSafeDto struct {
		ID           uint     `json:"id"`
		Dni          string   `json:"dni"`
		Mail         string   `json:"email"`
		Username     string   `json:"username"`
		CreatedRooms []r.Room `json:"createdRooms"`
	}
)
