package internal

import (
	"suffgo/database"
	r "suffgo/internal/room/entities"
	u "suffgo/internal/user/entities"
	ur "suffgo/internal/user_room/entities"
)

func Migrate(db database.Database) {
	err := db.GetDb().AutoMigrate(&u.User{}, &r.Room{}, &ur.UserRoom{})
	if err != nil {
		panic("Could migrate to database")
	}
}
