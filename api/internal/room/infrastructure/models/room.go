package models

type (
	Room struct {
		ID         uint    `xorm:"'id' pk autoincr"`
		LinkInvite *string `xorm:"'link_invite' varchar(255) null"`
		IsFormal   bool    `xorm:"not null"`
		Name       string  `xorm:"varchar(255) not null"`
		UserID     uint    `xorm:"'user_id' index not null"` //admin
	}
)
