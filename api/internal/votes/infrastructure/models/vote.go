package models

type Vote struct {
	ID       uint `xorm:"'id' pk autoincr"`
	UserID   uint `xorm:"'user_id' index not null"`
	OptionID uint `xorm:"'option_id' index not null"`
}
