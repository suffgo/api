package domain

import (
	v "suffgo/internal/rooms/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"
	"time"
)

type (
	Room struct {
		id          *sv.ID
		isFormal    v.IsFormal
		name        v.Name
		adminID     *sv.ID
		code        *v.InviteCode //es opcional porque al momento de creacion no existe
		description v.Description
		state       *v.State
		image       *v.Image
	}

	RoomDTO struct {
		ID          uint   `json:"id"`
		IsFormal    bool   `json:"is_formal"`
		Name        string `json:"name"`
		AdminID     uint   `json:"admin_id"`
		Description string `json:"description"`
		Code        string `json:"room_code"`
		State       string `json:"state"`
		Image       string `json:"image"`
	}

	//Dto para informacion util al frontend
	RoomDetailedDTO struct {
		ID          uint       `json:"id"`
		IsFormal    bool       `json:"is_formal"`
		RoomTitle   string     `json:"room_title"`
		AdminName   string     `json:"admin_name"`
		Description string     `json:"description"`
		Code        string     `json:"room_code"`
		State       string     `json:"state"`
		StartTime   *time.Time `json:"start_time"`
		Image       string     `json:"image"`
		Privileges  bool       `json:"privileges"`
	}

	RoomCreateRequest struct {
		IsFormal    bool   `json:"is_formal"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
	}

	RoomUpdate struct {
		IsFormal    bool   `json:"is_formal"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
	}

	JoinRoomRequest struct {
		RoomCode string `json:"room_code"`
	}

	AddSingleUserRequest struct {
		UserData string `json:"user_data"`
		RoomID   uint   `json:"room_id"`
	}

	RemoveFromWhitelistRequest struct {
		UserId uint `json:"user_id"`
		RoomId uint `json:"room_id"`
	}
)

func NewRoom(
	id *sv.ID,
	isFormal v.IsFormal,
	code *v.InviteCode,
	name v.Name,
	adminID *sv.ID,
	description v.Description,
	image *v.Image,
	state *v.State,
) *Room {
	return &Room{
		id:          id,
		isFormal:    isFormal,
		name:        name,
		adminID:     adminID,
		code:        code,
		description: description,
		image:       image,
		state:       state,
	}
}

func (r *Room) ID() sv.ID {
	return *r.id
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

func (r *Room) Code() v.InviteCode {

	return *r.code
}

func (r *Room) SetInviteCode(inviteCode v.InviteCode) {
	r.code = &inviteCode
}

func (r *Room) Description() v.Description {
	return r.description
}

func (r *Room) SetDescription(description v.Description) {
	r.description = description
}

func (r *Room) State() *v.State {
	return r.state
}

func (r *Room) Image() *v.Image {
	return r.image
}

func (r *Room) SetImage(image *v.Image) {
	r.image = image
}
