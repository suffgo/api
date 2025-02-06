package infrastructure

import (
	"fmt"
	"net/http"
	"suffgo/cmd/config"
	"suffgo/cmd/database"

	userUsecase "suffgo/internal/users/application/useCases"
	u "suffgo/internal/users/infrastructure"

	optionUsecase "suffgo/internal/options/application/useCases"
	o "suffgo/internal/options/infrastructure"

	voteUsecase "suffgo/internal/votes/application/useCases"
	v "suffgo/internal/votes/infrastructure"

	settingRoomUsecase "suffgo/internal/settingsRoom/application/useCases"
	sr "suffgo/internal/settingsRoom/infrastructure"

	proposalUsecase "suffgo/internal/proposals/application/useCases"
	p "suffgo/internal/proposals/infrastructure"

	roomUsecase "suffgo/internal/rooms/application/useCases"
	roomUsecaseAddUsers "suffgo/internal/rooms/application/useCases/addUsers"

	r "suffgo/internal/rooms/infrastructure"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

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
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:4321"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	// s.app.Pre(middleware.HTTPSNonWWWRedirect()) a tener en cuenta para el futuro en caso de despliegue

	authKey := []byte(s.conf.SecretKey)
	store := sessions.NewCookieStore(authKey)
	s.app.Use(session.Middleware(store))

	userRepo := s.InitializeUser()
	roomRepo := s.InitializeRoom(userRepo)
	s.InitializeProposal(roomRepo)
	s.InitializeVote()
	s.InitializeOption()
	s.InitializeSettingRoom()

	s.app.GET("/v1/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// for _, route := range s.app.Routes() {
	// 	fmt.Printf("Ruta registrada: Método=%s, Ruta=%s\n", route.Method, route.Path)
	// } estas lineas sirven para debuguear registros de rutas

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}

func (s *EchoServer) InitializeUser() *u.UserXormRepository {
	userRepo := u.NewUserXormRepository(s.db)

	// Initialize Use Cases
	createUserUseCase := userUsecase.NewCreateUsecase(userRepo)
	deleteUserUseCase := userUsecase.NewDeleteUsecase(userRepo)
	getAllUsersUseCase := userUsecase.NewGetAllUsecase(userRepo)
	getUserByIDUseCase := userUsecase.NewGetByIDUsecase(userRepo)
	getUserByEmail := userUsecase.NewGetByEmailUsecase(userRepo)
	loginUseCase := userUsecase.NewLoginUsecase(userRepo)
	restoreUseCase := userUsecase.NewRestoreUsecase(userRepo)
	changePasswordUseCase := userUsecase.NewChangePasswordUsecase(userRepo)
	updateUseCase := userUsecase.NewUpdateUsecase(userRepo)
	// Initialize Handler
	userHandler := u.NewUserEchoHandler(
		createUserUseCase,
		deleteUserUseCase,
		getAllUsersUseCase,
		getUserByIDUseCase,
		getUserByEmail,
		loginUseCase,
		restoreUseCase,
		changePasswordUseCase,
		updateUseCase,
	)

	// Initialize User Router
	u.InitializeUserEchoRouter(s.app, userHandler)

	return userRepo
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

func (s *EchoServer) InitializeRoom(userRepo *u.UserXormRepository) *r.RoomXormRepository {
	roomRepo := r.NewRoomXormRepository(s.db)
	createRoomUseCase := roomUsecase.NewCreateUsecase(roomRepo)
	deleteRoomUseCase := roomUsecase.NewDeleteUsecase(roomRepo)
	getAllRoomUseCase := roomUsecase.NewGetAllUsecase(roomRepo)
	getByIDRoomUseCase := roomUsecase.NewGetByIDUsecase(roomRepo)
	getByAdminRoomUseCase := roomUsecase.NewGetByAdminUsecase(roomRepo)
	restoreUseCase := roomUsecase.NewRestoreUsecase(roomRepo)
	joinUsecase := roomUsecase.NewJoinRoomUsecase(roomRepo)
	AddSingleUserUsecase := roomUsecaseAddUsers.NewAddSingleUserUsecase(roomRepo, userRepo)
	UpdateRoomUseCase := roomUsecase.NewUpdateRoomUsecase(roomRepo)

	roomHandler := r.NewRoomEchoHandler(
		createRoomUseCase,
		deleteRoomUseCase,
		getAllRoomUseCase,
		getByIDRoomUseCase,
		getByAdminRoomUseCase,
		restoreUseCase,
		joinUsecase,
		AddSingleUserUsecase,
		UpdateRoomUseCase,
	)
	r.InitializeRoomEchoRouter(s.app, roomHandler)

	return roomRepo

}

func (s *EchoServer) InitializeSettingRoom() {
	settingRoomRepo := sr.NewSettingRoomXormRepository(s.db)
	createSettingRoomUseCase := settingRoomUsecase.NewCreateUsecase(settingRoomRepo)
	deleteSettingRoomUseCase := settingRoomUsecase.NewDeleteUsecase(settingRoomRepo)
	getAllSettingRoomUseCase := settingRoomUsecase.NewGetAllUsecase(settingRoomRepo)
	getSettingRoomByIDUseCase := settingRoomUsecase.NewGetByIDUsecase(settingRoomRepo)
	updateSettingRoom := settingRoomUsecase.NewUpdateSettingRoomUsecase(settingRoomRepo)
	settingRoomHandler := sr.NewSettingRoomEchoHandler(
		createSettingRoomUseCase,
		deleteSettingRoomUseCase,
		getAllSettingRoomUseCase,
		getSettingRoomByIDUseCase,
		updateSettingRoom,
	)
	sr.InitializeSettingRoomEchoRouter(s.app, settingRoomHandler)
}

func (s *EchoServer) InitializeProposal(roomRepo *r.RoomXormRepository) {

	proposalRepo := p.NewProposalXormRepository(s.db)

	createProposalUseCase := proposalUsecase.NewCreateUsecase(proposalRepo, roomRepo)
	deleteProposalUseCase := proposalUsecase.NewDeleteUseCase(proposalRepo)
	getAllProposalsUseCase := proposalUsecase.NewGetAllUseCase(proposalRepo)
	getProposalByIDUseCase := proposalUsecase.NewGetByIDUseCase(proposalRepo)
	restoreProposalUseCase := proposalUsecase.NewRestoreUsecase(proposalRepo)

	proposalHandler := p.NewProposalEchoHandler(
		createProposalUseCase,
		getAllProposalsUseCase,
		getProposalByIDUseCase,
		deleteProposalUseCase,
		restoreProposalUseCase,
	)

	p.InitializeProposalEchoRouter(s.app, proposalHandler)
}
