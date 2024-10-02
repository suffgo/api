package main

import (
  "suffgo-backend-t/config"
  "suffgo-backend-t/database"
  server "suffgo-backend-t/server"
)

func main() {
  conf := config.GetConfig()
  db := database.NewPostgresDatabase(conf)
  server.NewEchoServer(conf, db).Start()
}