package domain

import (
	v "suffgo/internal/rooms/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	Room struct {
		id          *sv.ID
		linkInvite  v.LinkInvite
		isFormal    v.IsFormal
		name        v.Name
		adminID     *sv.ID
		inviteCode  *v.InviteCode //es opcional porque al momento de creacion no existe
		description v.Description
	}

	RoomDTO struct {
		ID          uint   `json:"id"`
		LinkInvite  string `json:"link_invite"`
		IsFormal    bool   `json:"is_formal"`
		Name        string `json:"name"`
		AdminID     uint   `json:"admin_id"`
		Description string `json:"description"`
		RoomCode    string `json:"room_code"`
	}

	RoomDetailedDTO struct {
		ID          uint   `json:"id"`
		LinkInvite  string `json:"link_invite"`
		RoomTitle   string `json:"room_title"` //es el nombre
		AdminName   string `json:"admin_name"`
		Description string `json:"description"`
		RoomCode    string `json:"room_code"`
	}

	RoomCreateRequest struct {
		LinkInvite  string `json:"link_invite"`
		IsFormal    bool   `json:"is_formal"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	JoinRoomRequest struct {
		RoomCode string `json:"room_code"`
	}

	AddSingleUserRequest struct {
		UserData string `json:"user_data"`
		RoomID   string `json:"room_id"`
	}
)

func NewRoom(
	id *sv.ID,
	linkInvite v.LinkInvite,
	isFormal v.IsFormal,
	name v.Name,
	adminID *sv.ID,
	description v.Description,
) *Room {
	return &Room{
		id:          id,
		linkInvite:  linkInvite,
		isFormal:    isFormal,
		name:        name,
		adminID:     adminID,
		inviteCode:  nil,
		description: description,
	}
}

func (r *Room) ID() sv.ID {
	return *r.id
}

func (r *Room) LinkInvite() v.LinkInvite {
	return r.linkInvite
}

func (r *Room) IsFormal() v.IsFormal {
	return r.isFormal
}

func (r *Room) SetIsFormal(isFormal v.IsFormal) {
	r.isFormal = isFormal
}

func (r *Room) Name() v.Name {
	return r.name
}

func (r *Room) SetName(name v.Name) {
	r.name = name
}

func (r *Room) AdminID() sv.ID {
	return *r.adminID
}

func (r *Room) InviteCode() v.InviteCode {

	return *r.inviteCode
}

func (r *Room) SetInviteCode(inviteCode v.InviteCode) {
	r.inviteCode = &inviteCode
}

func (r *Room) Description() v.Description {
	return r.description
}

func (r *Room) SetDescription(description v.Description) {
	r.description = description
}
