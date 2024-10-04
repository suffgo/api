package server

import (
	"fmt"
	"suffgo/config"
	"suffgo/database"

	userHandlers "suffgo/internal/user/handlers"
	userRepositories "suffgo/internal/user/repositories"
	userUsecases "suffgo/internal/user/usecases"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type echoServer struct {
	app  *echo.Echo
	db   database.Database
	conf *config.Config
}

func NewEchoServer(conf *config.Config, db database.Database) Server {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)

	return &echoServer{
		app:  echoApp,
		db:   db,
		conf: conf,
	}
}

func (s *echoServer) Start() {
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())

	// Health check adding
	s.app.GET("v1/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	s.initializeUserHttpHandler()
	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}

func (s *echoServer) initializeUserHttpHandler() {
	// Initialize all layers
	userPostgresRepository := userRepositories.NewUserPostgresRepository(s.db)
	userUsecase := userUsecases.NewUserUsecaseImpl(
		userPostgresRepository,
	)

	userHttpHandler := userHandlers.NewuserHttpHandler(userUsecase)

	// Routers
	userRoutes := s.app.Group("v1/user")
	userRoutes.POST("/register", userHttpHandler.RegisterUser)
	userRoutes.GET("/:id", userHttpHandler.GetUserByID)
	userRoutes.DELETE("/:id", userHttpHandler.DeleteUser)
	userRoutes.GET("", userHttpHandler.GetAll)
}
