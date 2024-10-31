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
		Server    *Server
		Db        *Db
		SecretKey string
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
		err := godotenv.Load("././.env")
		if err != nil {
			log.Fatalf("Error al cargar archivo .env: %v", err)
		}

		dbHost := os.Getenv("POSTGRES_HOST")
		dbUser := os.Getenv("POSTGRES_USER")
		dbPass := os.Getenv("POSTGRES_PASSWORD")
		dbName := os.Getenv("POSTGRES_DB")
		dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
		if err != nil {
			log.Fatalf("POSTGRES_PORT invalido: %s", os.Getenv("POSTGRES_PORT"))
		}

		apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
		if err != nil {
			log.Fatalf("API_PORT invalido: %s", os.Getenv("API_PORT"))
		}

		db := &Db{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPass,
			DBName:   dbName,
		}

		server := &Server{
			Port: apiPort,
		}

		configInstance = &Config{
			Server:    server,
			Db:        db,
		}
	})

	return configInstance
}


