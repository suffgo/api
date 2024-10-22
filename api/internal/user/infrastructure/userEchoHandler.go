package infrastructure

import (
	"net/http"
	u "suffgo/internal/user/application/useCases"

	"suffgo/internal/user/domain"
	v "suffgo/internal/user/domain/valueObjects"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	CreateUserUsecase  *u.CreateUsecase
	DeleteUserUsecase  *u.DeleteUsecase
	GetAllUsersUsecase *u.GetAllUsecase
	GetUserByIDUsecase *u.GetByIDUsecase
}

// Constructor for UserHandler
func NewUserEchoHandler(
	createUC *u.CreateUsecase,
	deleteUC *u.DeleteUsecase,
	getAllUC *u.GetAllUsecase,
	getByIDUC *u.GetByIDUsecase,
) *UserHandler {
	return &UserHandler{
		CreateUserUsecase:  createUC,
		DeleteUserUsecase:  deleteUC,
		GetAllUsersUsecase: getAllUC,
		GetUserByIDUsecase: getByIDUC,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	// Bind the request body to a DTO or model
	var req domain.UserCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// Map DTO to domain entity
	user := domain.NewUser(
		nil,
		*v.NewUserFullName(req.Name, req.Lastname),
		*v.NewUserUserName(req.Username),
		*v.NewUserDni(req.Dni),
		*v.NewUserEmail(req.Email),
		*v.NewUserPassword(req.Password),
	)

	// Call the use case
	err := h.CreateUserUsecase.Execute(*user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	return nil
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	return nil
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	return nil
}
