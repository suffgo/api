package infrastructure

import (
	"fmt"
	"suffgo/cmd/config"
	"suffgo/cmd/database"

	userUsecase "suffgo/internal/users/application/useCases"
	u "suffgo/internal/users/infrastructure"

	optionUsecase "suffgo/internal/options/application/useCases"
	o "suffgo/internal/options/infrastructure"

	voteUsecase "suffgo/internal/votes/application/useCases"
	v "suffgo/internal/votes/infrastructure"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	app  *echo.Echo
	db   database.Database
	conf *config.Config
}

var Db database.Database

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
	s.db.GetDb().ShowSQL(true)

	s.InitializeUser()
	s.InitializeOption()
	s.InitializeVote()

	// Health check adding
	s.app.GET("/v1/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// for _, route := range s.app.Routes() {
	// 	fmt.Printf("Ruta registrada: Método=%s, Ruta=%s\n", route.Method, route.Path)
	// }

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}

func (s *EchoServer) InitializeUser() {
	// Initialize the User Repository with xorm impl
	userRepo := u.NewUserXormRepository(s.db)

	// Le digo la base de datos a la que van a apuntar las operaciones hechas para jwt

	// Initialize Use Cases
	createUserUseCase := userUsecase.NewCreateUsecase(userRepo)
	deleteUserUseCase := userUsecase.NewDeleteUsecase(userRepo)
	getAllUsersUseCase := userUsecase.NewGetAllUsecase(userRepo)
	getUserByIDUseCase := userUsecase.NewGetByIDUsecase(userRepo)
	loginUseCase := userUsecase.NewLoginUsecase(userRepo)

	// Initialize Handler
	userHandler := u.NewUserEchoHandler(
		createUserUseCase,
		deleteUserUseCase,
		getAllUsersUseCase,
		getUserByIDUseCase,
		loginUseCase,
	)

	// Initialize User Router
	u.InitializeUserEchoRouter(s.app, userHandler)
}

func (s *EchoServer) InitializeOption() {
	optionRepo := o.NewOptionXormRepository(s.db)

	createOptionUseCase := optionUsecase.NewCreateUsecase(optionRepo)
	deleteOptionUseCase := optionUsecase.NewDeleteUsecase(optionRepo)
	getAllOptionUseCase := optionUsecase.NewGetAllRepository(optionRepo)
	getOptionByIDUseCase := optionUsecase.NewGetByIDUsecase(optionRepo)
	getOptionByValueUseCase := optionUsecase.NewGetByValueUsecase(optionRepo)

	optionHandler := o.NewOptionEchoHandler(
		createOptionUseCase,
		deleteOptionUseCase,
		getAllOptionUseCase,
		getOptionByIDUseCase,
		getOptionByValueUseCase,
	)
	o.InitializeOptionEchoRouter(s.app, optionHandler)
}

func (s *EchoServer) InitializeVote() {
	voteRepo := v.NewVoteXormRepository(s.db)

	createVoteUseCase := voteUsecase.NewCreateUsecase(voteRepo)
	deleteVoteUseCase := voteUsecase.NewDeleteUsecase(voteRepo)
	getAllVoteUseCase := voteUsecase.NewGetAllRepository(voteRepo)
	getVoteByIDUseCase := voteUsecase.NewGetByIDUsecase(voteRepo)

	voteHandler := v.NewVoteEchoHandler(
		createVoteUseCase,
		deleteVoteUseCase,
		getAllVoteUseCase,
		getVoteByIDUseCase,
	)
	v.InitializeVoteEchoRouter(s.app, voteHandler)
}
