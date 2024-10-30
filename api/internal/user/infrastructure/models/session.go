package models

type Session struct {
	ID        uint   `xorm:"'id' pk autoincr"`
	SessionID string `xorm:"'session_id' not null"`
	UserID    uint   `xorm:"'user_id' index not null"`
}