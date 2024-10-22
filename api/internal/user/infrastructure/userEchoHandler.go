package infrastructure

import (
	"fmt"
	"net/http"
	"strconv"
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

	fullname, err := v.NewFullName(req.Name, req.Lastname)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	username, err := v.NewUserName(req.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	dni, err := v.NewDni(req.Dni)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	email, err := v.NewEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	password, err := v.NewPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// Map DTO to domain entity
	user := domain.NewUser(
		nil,
		*fullname,
		*username,
		*dni,
		*email,
		*password,
	)

	// Call the use case
	err = h.CreateUserUsecase.Execute(*user)
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
	idParam := c.Param("id")
    idInput, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
    }
	
	id, _ := v.NewID(uint(idInput))
	user, err := h.GetUserByIDUsecase.Execute(*id)
	
	fmt.Printf("id = %d\n", id.Id)
	if err != nil {
        if err.Error() == "user not found" {
            return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User fetch error"})
    }

	userDTO := &domain.UserDTO{
		ID: user.ID().Id,
		Name: user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Username: user.Username().Username,
		Dni: user.Dni().Dni,
		Email: user.Email().Email,
		Password: user.Password().Password,
	}
	return c.JSON(http.StatusOK, userDTO)
}
