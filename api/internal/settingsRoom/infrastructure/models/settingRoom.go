package models

import "time"

type SettingRoom struct {
	ID            uint       `xorm:"'id' pk autoincr"`
	Quorum        *int       `xorm:"'quorum' null"`
	Privacy       bool       `xorm:"'privacy' not null default false"`
	VoterLimit    int        `xorm:"'voter_limit' not null default 0"`
	TimeDate      *time.Time `xorm:"'time_date' null"`
	ProposalTimer int        `xorm:"'proposal_timer' not null default 60"` //despues vemos que onda si es minutos o segundos
	RoomID        uint       `xorm:"'room_id' index not null"`
}