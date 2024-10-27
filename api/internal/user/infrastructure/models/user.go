package models

type Users struct {
	ID            uint   `xorm:"'id' pk autoincr"`
	Dni           string `xorm:"varchar(10) not null unique"`
	Username      string `xorm:"'username' varchar(50) not null unique"`
	Password      string `xorm:"varchar(255) not null"`
	Name          string `xorm:"varchar(255) not null"`
	Lastname      string `xorm:"'last_name' varchar(255) not null"`
	Email         string `xorm:"varchar(255) not null unique"`
}