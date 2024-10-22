package domain

import (
	v "suffgo/internal/user/domain/valueObjects"
)

type (
	User struct {
		id       v.UserID
		name     v.UserFullName
		username v.UserUserName
		dni      v.UserDni
		email    v.UserEmail
		password v.UserPassword
	}

	UserDTO struct {
		ID      v.UserID       `json:"id"`
		Name     v.UserFullName `json:"name"`
		Username v.UserUserName `json:"username"`
		Dni      v.UserDni      `json:"dni"`
		Email    v.UserEmail    `json:"email"`
		Password v.UserPassword `json:"password"`
	}
)

func NewUser(id v.UserID, name v.UserFullName, username v.UserUserName, dni v.UserDni, email v.UserEmail, password v.UserPassword) *User {
	return &User{
		id:       id,
		name:     name,
		username: username,
		dni:      dni,
		email:    email,
		password: password,
	}
}

func (u *User) ID() v.UserID {
	return u.id
}

func (u *User) Email() v.UserEmail {
	return u.email
}

func (u *User) Username() v.UserUserName {
	return u.username
}

func (u *User) Dni() v.UserDni {
	return u.dni
}

func (u *User) Password() v.UserPassword {
	return u.password
}

func (u *User) FullName() v.UserFullName {
	return u.name
}


