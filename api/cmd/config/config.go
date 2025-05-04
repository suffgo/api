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
		Prod      bool
	}

	Server struct {
		Port        int
		AllowedCORS string
		UploadsDir  string
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

		//secret key para validar sesiones
		secretKey := os.Getenv("SECRET_SESSION_AUTH_KEY")

		db := &Db{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPass,
			DBName:   dbName,
		}

		origins := os.Getenv("ALLOWED_CORS")

		server := &Server{
			Port:        apiPort,
			AllowedCORS: origins,
			UploadsDir:  os.Getenv("UPLOADS_DIR"),
		}

		configInstance = &Config{
			Server:    server,
			Db:        db,
			SecretKey: secretKey,
			Prod:      os.Getenv("PROD") == "true",
		}
	})

	return configInstance
}
