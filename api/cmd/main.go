package main

import (
	"suffgo/cmd/config"
	"suffgo/cmd/database"
	e "suffgo/internal/shared/infraestructure"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	e.NewEchoServer(db, conf).Start()
}
