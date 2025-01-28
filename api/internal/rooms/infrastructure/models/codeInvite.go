package models

// Definion de modelo de codigo de invitacion para las salas/rooms

type InviteCode struct {
	ID     uint `xorm:"'id' pk autoincr"`
	RoomID uint `xorm:"'room_id' index not null"`
	Code   string `xorm:"'code' unique not null varchar(255)"`
}
