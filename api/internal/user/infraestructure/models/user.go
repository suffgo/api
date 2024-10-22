package models

type User struct {
	ID       uint   `xorm:"pk autoincr'id'"`
	Dni      string `xorm:"varchar(10) notnull unique"`
	Username string `xorm:"varchar(50) notnull unique"`
	Password string `xorm:"varchar(255) notnull unique"`
	Name     string `xorm:"varchar(255) notnull"`
	Lastname string `xorm:"varchar(255) notnull"`
	Email    string `xorm:"varchar(255) notnull unique"`
}