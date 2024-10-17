package main

import (
	"suffgo/config"
	"suffgo/database"
	"suffgo/internal/user_room/entities"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	userMigrate(db)
}

func userMigrate(db database.Database) {
	db.GetDb().Migrator().CreateTable(&entities.UserRoom{})
}
