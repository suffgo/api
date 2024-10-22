package domain

import (
	v "suffgo/internal/user/domain/valueObjects"
)

type UserRepository interface {
	GetByID(id v.ID) (*User, error)
	GetAll() ([]User, error)
	Delete(id v.ID) error
	Create(user User) error
	GetByEmail(email v.Email) (*User, error)
	Save(user User) error
	// Update(user User) error
}
