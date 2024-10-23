package infrastructure

import (
	"fmt"
	"suffgo/cmd/config"
	"suffgo/cmd/database"

	userUsecase "suffgo/internal/user/application/useCases"
	u "suffgo/internal/user/infrastructure"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	app  *echo.Echo
	db   database.Database
	conf *config.Config
}

func NewEchoServer(db database.Database, conf *config.Config) *EchoServer {
	echoApp := echo.New()
	return &EchoServer{
		app:  echoApp,
		db:   db,
		conf: conf,
	}
}

func (s *EchoServer) Start() {
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())
	
	// Initialize the User Repository with xorm impl
	userRepo := u.NewUserXormRepository(s.db)

	// Initialize Use Cases
	createUserUseCase := userUsecase.NewCreateUsecase(userRepo)
	deleteUserUseCase := userUsecase.NewDeleteUsecase(userRepo)
	getAllUsersUseCase := userUsecase.NewGetAllUsecase(userRepo)
	getUserByIDUseCase := userUsecase.NewGetByIDUsecase(userRepo)

	// Initialize Handler
	userHandler := u.NewUserEchoHandler(
		createUserUseCase,
		deleteUserUseCase,
		getAllUsersUseCase,
		getUserByIDUseCase,
	)

	// Initialize User Router
	u.InitializeUserEchoRouter(s.app, userHandler)

	// Health check adding
	s.app.GET("/v1/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	for _, route := range s.app.Routes() {
		fmt.Printf("Ruta registrada: Método=%s, Ruta=%s\n", route.Method, route.Path)
	}
	
	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}
