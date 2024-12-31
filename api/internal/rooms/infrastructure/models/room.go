package models

import "time"

type (
	Room struct {
		ID         uint       `xorm:"'id' pk autoincr"`
		LinkInvite *string    `xorm:"'link_invite' varchar(255) null"`
		IsFormal   bool       `xorm:"'is_formal' not null"`
		Name       string     `xorm:"'name' varchar(255) not null"`
		AdminID    uint       `xorm:"'admin_id' index not null"` //admin
		DeleteAT   *time.Time `xorm:"deleted"`
	}
)
