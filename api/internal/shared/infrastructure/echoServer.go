package infrastructure

import (
	"fmt"
	"net/http"
	"strings"
	"suffgo/cmd/config"
	"suffgo/cmd/database"

	optDom "suffgo/internal/options/domain"
	propDom "suffgo/internal/proposals/domain"
	roomDom "suffgo/internal/rooms/domain"
	srDom "suffgo/internal/settingsRoom/domain"
	userDom "suffgo/internal/users/domain"
	voteDom "suffgo/internal/votes/domain"

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
	roomWsUsecase "suffgo/internal/rooms/application/useCases/websocket"

	r "suffgo/internal/rooms/infrastructure"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"suffgo/cmd/migrateFunc"
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

type Dependencies struct {
	UserRepo        userDom.UserRepository
	RoomRepo        roomDom.RoomRepository
	SettingRoomRepo srDom.SettingRoomRepository
	ProposalRepo    propDom.ProposalRepository
	VotesRepo       voteDom.VoteRepository
	OptionsRepo     optDom.OptionRepository
}

func NewDependencies(db database.Database) *Dependencies {
	userRepo := u.NewUserXormRepository(db)
	roomRepo := r.NewRoomXormRepository(db)
	settingRoomRepo := sr.NewSettingRoomXormRepository(db)
	proposalRepo := p.NewProposalXormRepository(db)
	voteRepo := v.NewVoteXormRepository(db)
	optionRepo := o.NewOptionXormRepository(db)

	return &Dependencies{
		UserRepo:        userRepo,
		RoomRepo:        roomRepo,
		SettingRoomRepo: settingRoomRepo,
		ProposalRepo:    proposalRepo,
		VotesRepo:       voteRepo,
		OptionsRepo:     optionRepo,
	}
}

func (s *EchoServer) Start() {

	if s.conf.Prod {
		origins := strings.Split(s.conf.Server.AllowedCORS, ",")

		s.app.Pre(middleware.HTTPSNonWWWRedirect()) //para redirigir http:// → https:// y eliminar www
		s.app.Pre(middleware.RemoveTrailingSlash())

		s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     origins,
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderCookie},
			ExposeHeaders:    []string{"Set-Cookie"},
			AllowCredentials: true,
		}))

		s.app.Static("/uploads", "internal/uploads/")

		if err := migrateFunc.Make(); err != nil {
			fmt.Printf("Migraciones ya fueron hechas: %v\n", err)
		}

		authKey := []byte(s.conf.SecretKey)
		store := sessions.NewCookieStore(authKey)
		store.Options = &sessions.Options{
			HttpOnly: true,
			Secure:   s.conf.Prod,
			SameSite: http.SameSiteNoneMode,
			Path:     "/",
		}
	
		s.app.Use(session.Middleware(store))

	} else {
		s.app.Debug = true
		s.db.GetDb().ShowSQL(true)
		s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:4321"},
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowCredentials: true,
		}))

		s.app.Static("/uploads", "internal/uploads/")

		authKey := []byte(s.conf.SecretKey)
		store := sessions.NewCookieStore(authKey)
		store.Options = &sessions.Options{
			HttpOnly: true,
			Secure:   s.conf.Prod,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		}
	
		s.app.Use(session.Middleware(store))
	}


	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())


	deps := NewDependencies(s.db)

	s.InitializeUser(deps.UserRepo, deps.RoomRepo, deps.SettingRoomRepo)
	s.InitializeRoom(deps.UserRepo, deps.SettingRoomRepo, deps.ProposalRepo, deps.OptionsRepo, deps.VotesRepo)
	s.InitializeSettingRoom(deps.SettingRoomRepo, deps.RoomRepo)
	s.InitializeProposal(deps.ProposalRepo, deps.RoomRepo)
	s.InitializeVote()
	s.InitializeOption()

	s.app.GET("/v1/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}

var getUserByIDUseCase *userUsecase.GetByIDUsecase

func (s *EchoServer) InitializeUser(userRepo userDom.UserRepository, roomRepo roomDom.RoomRepository, setrRepo srDom.SettingRoomRepository) {

	// Initialize Use Cases
	createUserUseCase := userUsecase.NewCreateUsecase(userRepo)
	deleteUserUseCase := userUsecase.NewDeleteUsecase(userRepo)
	getAllUsersUseCase := userUsecase.NewGetAllUsecase(userRepo)
	getUserByEmail := userUsecase.NewGetByEmailUsecase(userRepo)
	getUserByIDUseCase = userUsecase.NewGetByIDUsecase(userRepo)
	loginUseCase := userUsecase.NewLoginUsecase(userRepo)
	restoreUseCase := userUsecase.NewRestoreUsecase(userRepo)
	changePasswordUseCase := userUsecase.NewChangePasswordUsecase(userRepo)
	updateUseCase := userUsecase.NewUpdateUsecase(userRepo)
	getByRoom := userUsecase.NewGetUsersByRoom(userRepo, roomRepo, setrRepo)
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
		getByRoom,
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
	getOptionByPropUsecase := optionUsecase.NewGetByPropUsecase(optionRepo)

	optionHandler := o.NewOptionEchoHandler(
		createOptionUseCase,
		deleteOptionUseCase,
		getAllOptionUseCase,
		getOptionByIDUseCase,
		getOptionByValueUseCase,
		getOptionByPropUsecase,
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

func (s *EchoServer) InitializeRoom(
	userRepo userDom.UserRepository,
	settingRoomRepo srDom.SettingRoomRepository,
	proposalRepo propDom.ProposalRepository,
	optionsRepo optDom.OptionRepository,
	votesRepo voteDom.VoteRepository,
) {
	roomRepo := r.NewRoomXormRepository(s.db)
	createRoomUC := roomUsecase.NewCreateUsecase(roomRepo, settingRoomRepo)
	deleteRoomUC := roomUsecase.NewDeleteUsecase(roomRepo)
	getAllRoomUC := roomUsecase.NewGetAllUsecase(roomRepo)
	getByIDRoomUC := roomUsecase.NewGetByIDUsecase(roomRepo)
	getByAdminRoomUC := roomUsecase.NewGetByAdminUsecase(roomRepo)
	restoreUC := roomUsecase.NewRestoreUsecase(roomRepo)
	joinUC := roomUsecase.NewJoinRoomUsecase(roomRepo, settingRoomRepo)
	AddSingleUserUC := roomUsecaseAddUsers.NewAddSingleUserUsecase(roomRepo, userRepo)
	UpdateRoomUC := roomUsecase.NewUpdateRoomUsecase(roomRepo)
	ManageWsUC := roomWsUsecase.NewManageWsUsecase(roomRepo, userRepo, proposalRepo, optionsRepo, votesRepo)
	getSrByRoomIDUC := roomUsecase.NewGetSrByRoomUsecase(roomRepo, settingRoomRepo)
	HistoryUC := roomUsecase.NewHistoryRoomsUsecase(roomRepo)
	rmWhitelistUC := roomUsecase.NewWhitelistRmUsecase(roomRepo, userRepo)

	roomHandler := r.NewRoomEchoHandler(
		createRoomUC,
		deleteRoomUC,
		getAllRoomUC,
		getByIDRoomUC,
		getByAdminRoomUC,
		restoreUC,
		joinUC,
		AddSingleUserUC,
		getUserByIDUseCase,
		UpdateRoomUC,
		ManageWsUC,
		getSrByRoomIDUC,
		rmWhitelistUC,
		HistoryUC,
	)
	r.InitializeRoomEchoRouter(s.app, roomHandler)

}

func (s *EchoServer) InitializeSettingRoom(srRepo srDom.SettingRoomRepository, roomRepo roomDom.RoomRepository) {

	createSettingRoomUseCase := settingRoomUsecase.NewCreateUsecase(srRepo, roomRepo)
	deleteSettingRoomUseCase := settingRoomUsecase.NewDeleteUsecase(srRepo)
	getAllSettingRoomUseCase := settingRoomUsecase.NewGetAllUsecase(srRepo)
	getSettingRoomByIDUseCase := settingRoomUsecase.NewGetByIDUsecase(srRepo)
	updateSettingRoom := settingRoomUsecase.NewUpdateSettingRoomUsecase(srRepo, roomRepo)
	getByRoomIdUsecase := settingRoomUsecase.NewGetByRoomID(srRepo)
	settingRoomHandler := sr.NewSettingRoomEchoHandler(
		createSettingRoomUseCase,
		deleteSettingRoomUseCase,
		getAllSettingRoomUseCase,
		getSettingRoomByIDUseCase,
		updateSettingRoom,
		getByRoomIdUsecase,
	)
	sr.InitializeSettingRoomEchoRouter(s.app, settingRoomHandler)
}

func (s *EchoServer) InitializeProposal(propRepo propDom.ProposalRepository, roomRepo roomDom.RoomRepository) {

	createProposalUseCase := proposalUsecase.NewCreateUsecase(propRepo, roomRepo)
	deleteProposalUseCase := proposalUsecase.NewDeleteUseCase(propRepo, roomRepo)
	getAllProposalsUseCase := proposalUsecase.NewGetAllUseCase(propRepo)
	getProposalByIDUseCase := proposalUsecase.NewGetByIDUseCase(propRepo)
	updateProposalUseCase := proposalUsecase.NewUpdateProposalUsecase(propRepo, roomRepo)
	getByRoomUsecase := proposalUsecase.NewGetByRoomUsecase(propRepo)
	getResultsByRoomUsecase := proposalUsecase.NewGetResultsByRoomUsecase(propRepo)

	proposalHandler := p.NewProposalEchoHandler(
		createProposalUseCase,
		getAllProposalsUseCase,
		getProposalByIDUseCase,
		deleteProposalUseCase,
		updateProposalUseCase,
		getByRoomUsecase,
		getResultsByRoomUsecase,
	)

	p.InitializeProposalEchoRouter(s.app, proposalHandler)
}
