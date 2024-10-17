package main

import (
	"suffgo/config"
	"suffgo/database"
	"suffgo/internal"
	"suffgo/server"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	internal.Migrate(db)
	server.NewEchoServer(conf, db).Start()
}
