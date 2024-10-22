package database

import (
	"fmt"
	"suffgo/cmd/config"
	"sync"

	"xorm.io/xorm"
	_ "github.com/lib/pq"
)

type postgresDatabase struct {
	Db *xorm.Engine
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

var engine *xorm.Engine

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/Argentina/Buenos_Aires",
			conf.Db.Host,
			conf.Db.User,
			conf.Db.Password,
			conf.Db.DBName,
		)

		engine, err := xorm.NewEngine("postgres", dsn)
		if err != nil {
			panic("error to connect database")
		}

		dbInstance = &postgresDatabase{Db: engine}
	})

	return dbInstance
}

func (p *postgresDatabase) GetDb() *xorm.Engine {
	return dbInstance.Db
}
