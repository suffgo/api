package infraestructure

import (
	"net/http"
	u "suffgo/internal/user/application/useCases"

	"github.com/labstack/echo/v4"
	"suffgo/internal/user/domain"
	v "suffgo/internal/user/domain/valueObjects"
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
	var userDTO domain.UserDTO
	if err := c.Bind(&userDTO); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Map DTO to domain entity
	user := domain.NewUser(
		*v.NewUserID(userDTO.ID.Id),
		*v.NewUserFullName(userDTO.Name.Name, userDTO.Name.Lastname ),
		*v.NewUserUserName(userDTO.Username.Username),
		*v.NewUserDni(userDTO.Dni.Dni),
		*v.NewUserEmail(userDTO.Email.Email),
		*v.NewUserPassword(userDTO.Password.Password),
	)

	// Call the use case
	err := h.CreateUserUsecase.Execute(*user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, userDTO)
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