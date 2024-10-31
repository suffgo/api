package models

type UserRoom struct {
	ID     uint `xorm:"'id' pk autoincr"`
	UserID uint `xorm:"'user_id' index not null"` // Usuario habilitado
	RoomID uint `xorm:"'room_id' index not null"` // Sala habilitada
}

