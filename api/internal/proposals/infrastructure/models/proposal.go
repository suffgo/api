package models

import "time"

type Proposal struct {
	ID          uint       `xorm:"'id' pk autoincr"`
	Archive     *string    `xorm:"'archive' null"` // Archivo con informacion detallada de la propuesta
	Title       string     `xorm:"'title' not null"`
	Description *string    `xorm:"'description' null"`
	RoomID      uint       `xorm:"'room_id' index not null"`
	DeleteAT    *time.Time `xorm:"deleted"`
}
