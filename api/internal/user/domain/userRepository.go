package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	v "suffgo/internal/user/domain/valueObjects"
)

type UserRepository interface {
	GetByID(id sv.ID) (*User, error)
	GetAll() ([]User, error)
	Delete(id sv.ID) error
	GetByEmail(email v.Email) (*User, error)
	Save(user User) error
	GetByDni(dni v.Dni) (*User, error)
	GetByUsername(username v.UserName) (*User, error)
	GetIDBySession(session string, id sv.ID) (*sv.ID, error)
	SaveSession(session string, id sv.ID) error
	// Update(user User) error
}
