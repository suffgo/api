package domain

import (
	v "suffgo/internal/user/domain/valueObjects"
)

type UserRepository interface {
	GetByID(id v.UserID) (*User, error)
	GetAll() ([]User, error)
	Delete(id v.UserID) error
	Create(user User) error
	GetByEmail(email v.UserEmail) (*User, error)
	Save(user User) error
	// Update(user User) error
}