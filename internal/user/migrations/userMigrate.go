package main

import (
	"suffgo-backend-t/config"
	"suffgo-backend-t/database"
	"suffgo-backend-t/internal/user/entities"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	userMigrate(db)
}

func userMigrate(db database.Database) {
	db.GetDb().Migrator().CreateTable(&entities.User{})
}
