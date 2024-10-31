package domain

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	v "suffgo/internal/users/domain/valueObjects"
)

type (
	User struct {
		id       *sv.ID
		name     v.FullName
		username v.UserName
		dni      v.Dni
		email    v.Email
		password v.Password
	}

	UserDTO struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Username string `json:"username"`
		Dni      string `json:"dni"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UserCreateRequest struct {
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Username string `json:"username"`
		Dni      string `json:"dni"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

func NewUser(
	id *sv.ID,
	name v.FullName,
	username v.UserName,
	dni v.Dni,
	email v.Email,
	password v.Password,
) *User {
	return &User{
		id:       id,
		name:     name,
		username: username,
		dni:      dni,
		email:    email,
		password: password,
	}
}

func (u *User) ID() sv.ID {
	return *u.id
}

func (u *User) Email() v.Email {
	return u.email
}

func (u *User) Username() v.UserName {
	return u.username
}

func (u *User) Dni() v.Dni {
	return u.dni
}

func (u *User) Password() v.Password {
	return u.password
}

func (u *User) FullName() v.FullName {
	return u.name
}
