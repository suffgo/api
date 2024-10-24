package domain

import (
	v "suffgo/internal/user/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type UserRepository interface {
	GetByID(id sv.ID) (*User, error)
	GetAll() ([]User, error)
	Delete(id sv.ID) error
	GetByEmail(email v.Email) (*User, error)
	Save(user User) error
	// Update(user User) error
}
