package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Server *Server
		Db     *Db
	}

	Server struct {
		Port int
	}

	Db struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {
	once.Do(func() {
      err := godotenv.Load()
      if err != nil {
         log.Fatalf("Error al cargar archivo .env: %v", err)
      }

      dbHost := os.Getenv("DB_HOST")
      dbUser := os.Getenv("DB_USER")
      dbPass := os.Getenv("DB_PASS")
      dbName := os.Getenv("DB_NAME")
      dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
      
      if err != nil {
         log.Fatalf("DB_PORT invalido: %s", os.Getenv("DB_PORT"))
      }


      apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))

      if err != nil {
         log.Fatalf("API_PORT invalido: %s", os.Getenv("API_PORT"))
      }

      db := &Db{
         Host: dbHost,
         Port: dbPort,
         User: dbUser,
         Password: dbPass,
         DBName: dbName,
      }

      server := &Server{
         Port: apiPort,
      }

      configInstance = &Config{
         Server: server,
         Db: db,
      } 
	})

	return configInstance
}
