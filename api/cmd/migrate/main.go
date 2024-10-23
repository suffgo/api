package main

import (
	"fmt"
	"suffgo/cmd/config"
	"suffgo/cmd/database"
	m "suffgo/internal/user/infrastructure/models"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)

	MigrateUser(db)
	//MigrateRoom(db)
	//MigrateProposal(db)
	//MigrateOption(db)
	//MigrateRoomSetting(db)
	//MigrateElection(db)
}

func MigrateUser(db database.Database) error {
	err := db.GetDb().Sync2(new(m.User))

	if err!= nil {
		panic(err)
	} else {
		fmt.Printf("Se ha migrado User con exito\n")
	}

	return err
}