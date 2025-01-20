package domain

import (
	v "suffgo/internal/rooms/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	Room struct {
		id         *sv.ID
		linkInvite v.LinkInvite
		isFormal   v.IsFormal
		name       v.Name
		adminID    *sv.ID
		inviteCode *v.InviteCode //es opcional porque al momento de creacion no existe
	}

	RoomDTO struct {
		ID         uint   `json:"id"`
		LinkInvite string `json:"link_invite"`
		IsFormal   bool   `json:"is_formal"`
		Name       string `json:"name"`
		AdminID    uint   `json:"admin_id"`
		RoomCode string `json:"room_code"`
	}

	RoomCreateRequest struct {
		LinkInvite string `json:"link_invite"`
		IsFormal   bool   `json:"is_formal"`
		Name       string `json:"name"`
	}

	JoinRoomRequest struct {
		RoomCode string `json:"room_code"`
	}
)

func NewRoom(
	id *sv.ID,
	linkInvite v.LinkInvite,
	isFormal v.IsFormal,
	name v.Name,
	adminID *sv.ID,
) *Room {
	return &Room{
		id:         id,
		linkInvite: linkInvite,
		isFormal:   isFormal,
		name:       name,
		adminID:    adminID,
		inviteCode: nil,
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

func (r *Room) Name() v.Name {
	return r.name
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