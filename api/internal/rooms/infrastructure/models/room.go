package models

import "time"

type (
	Room struct {
		ID          uint       `xorm:"'id' pk autoincr"`
		LinkInvite  *string    `xorm:"'link_invite' varchar(255) null"`
		IsFormal    bool       `xorm:"'is_formal' not null"`
		Name        string     `xorm:"'name' varchar(255) not null"`
		AdminID     uint       `xorm:"'admin_id' index not null"` //admin
		Description string     `xorm:"'description' varchar(255) not null"`
		State       string     `xorm:"'state' varchar(16) not null"`
		Image       string     `xorm:"'image' varchar null"`
		DeletedAt   *time.Time `xorm:"deleted"`
	}
)
