package models

import "time"

type SettingsRoom struct {
	ID            uint       `xorm:"'id' pk autoincr"`
	Quorum        *int       `xorm:"'quorum' null"`
	Privacy       bool       `xorm:"'privacy' not null default false"`
	VoterLimit    int        `xorm:"'voter_limit' not null default 0"`
	StartTime     *time.Time `xorm:"'start_time' null"`
	ProposalTimer int        `xorm:"'proposal_timer' not null default 60"` //despues vemos que onda si es minutos o segundos
	RoomID        uint       `xorm:"'room_id' index not null"`
}
