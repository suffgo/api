package database

import (
	"fmt"
	"log"
	"suffgo/cmd/config"
	"sync"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
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
        // Construyes el DSN  
        dsn := fmt.Sprintf(
            "host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/Argentina/Buenos_Aires",
            conf.Db.Host,
            conf.Db.User,
            conf.Db.Password,
            conf.Db.DBName,
        )

        // 🔍 Loggea las variables críticas
        log.Printf("▶️  POSTGRES HOST    = %q", conf.Db.Host)
        log.Printf("▶️  POSTGRES USER    = %q", conf.Db.User)
        log.Printf("▶️  POSTGRES PASS    = %q", mask(conf.Db.Password))
        log.Printf("▶️  POSTGRES DBNAME  = %q", conf.Db.DBName)
        log.Printf("▶️  DSN              = %q", dsn)

        engine, err := xorm.NewEngine("postgres", dsn)
        if err != nil {
            log.Fatalf("❌ error to connect database: %v", err)
        }

        dbInstance = &postgresDatabase{Db: engine}
    })

    return dbInstance
}

func mask(s string) any {
	panic("unimplemented")
}

func (p *postgresDatabase) GetDb() *xorm.Engine {
	return dbInstance.Db
}
