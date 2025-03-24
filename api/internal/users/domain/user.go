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
		image    *v.Image
	}

	UserDTO struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Username string `json:"username"`
		Dni      string `json:"dni"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Image    string `json:"image"`
	}

	//para no comprometer la pass
	UserSafeDTO struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Username string `json:"username"`
		Dni      string `json:"dni"`
		Email    string `json:"email"`
		Image    string `json:"image"`
	}

	UserCreateRequest struct {
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Username string `json:"username"`
		Dni      string `json:"dni"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Image    string `json:"image"`
	}

	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	ChangePasswordRequest struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}
)

func NewUser(
	id *sv.ID,
	name v.FullName,
	username v.UserName,
	dni v.Dni,
	email v.Email,
	password v.Password,
	image *v.Image,
) *User {
	return &User{
		id:       id,
		name:     name,
		username: username,
		dni:      dni,
		email:    email,
		password: password,
		image:    image,
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

func (u *User) Image() *v.Image {
	return u.image
}
