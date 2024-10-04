package main

import (
	"suffgo/config"
	"suffgo/database"
	server "suffgo/server"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	server.NewEchoServer(conf, db).Start()
}
