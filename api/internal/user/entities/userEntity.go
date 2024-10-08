package entities

import "gorm.io/gorm"

type (
	User struct {
		gorm.Model
		Dni      uint32 `gorm:"unique; not null"`
		Mail     string `gorm:"unique; not null"`
		Password string `gorm:"not null"`
		Username string `gorm:"unique; not null"`
	}

	UserDto struct {
		ID       uint64 `json:"id"`
		Dni      uint32
		Mail     string
		Password string
		Username string
	}

	UserSafeDto struct {
		ID       uint   `json:"id"`
		Dni      uint32 `json:"dni"`
		Mail     string `json:"email"`
		Username string `json:"username"`
	}
)
