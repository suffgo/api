package models

type Users struct {
	ID            uint   `xorm:"'id' pk autoincr"`
	Dni           string `xorm:"varchar(10) notnull unique"`
	Username      string `xorm:"'username' varchar(50) notnull unique"`
	Password      string `xorm:"varchar(255) notnull"`
	Name          string `xorm:"varchar(255) notnull"`
	Lastname      string `xorm:"'last_name' varchar(255) notnull"`
	Email         string `xorm:"varchar(255) notnull unique"`
}
