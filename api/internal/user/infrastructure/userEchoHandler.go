package infrastructure

import (
	"fmt"
	"net/http"
	"strconv"
	u "suffgo/internal/user/application/useCases"

	d "suffgo/internal/user/domain"
	v "suffgo/internal/user/domain/valueObjects"

	sv "suffgo/internal/shared/domain/valueObjects"

	se "suffgo/internal/shared/domain/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserEchoHandler struct {
	CreateUserUsecase  *u.CreateUsecase
	DeleteUserUsecase  *u.DeleteUsecase
	GetAllUsersUsecase *u.GetAllUsecase
	GetUserByIDUsecase *u.GetByIDUsecase
	LoginUsecase       *u.LoginUsecase
	ValidateUsecase    *u.ValidateSessionUsecase
}

// Constructor for UserEchoHandler
func NewUserEchoHandler(
	createUC *u.CreateUsecase,
	deleteUC *u.DeleteUsecase,
	getAllUC *u.GetAllUsecase,
	getByIDUC *u.GetByIDUsecase,
	loginUC *u.LoginUsecase,
	ValidateUC *u.ValidateSessionUsecase,
) *UserEchoHandler {
	return &UserEchoHandler{
		CreateUserUsecase:  createUC,
		DeleteUserUsecase:  deleteUC,
		GetAllUsersUsecase: getAllUC,
		GetUserByIDUsecase: getByIDUC,
		LoginUsecase:       loginUC,
		ValidateUsecase:    ValidateUC,
	}
}

func (u *UserEchoHandler) Login(c echo.Context) error {

	var req d.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	username, err := v.NewUserName(req.Username)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	pass, err := v.NewPassword(req.Password)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	userID, err := u.LoginUsecase.Execute(*username, *pass)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}

	token, err := createToken(*username, c.RealIP(), c.Request().UserAgent(), *userID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// Devolver el token al cliente
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (h *UserEchoHandler) SecureHello(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)

	// Extraer los claims del token
	claims := user.Claims.(jwt.MapClaims)
	name := claims["username"].(string)

	if claims["ip"] != c.RealIP() || claims["user_agent"] != c.Request().UserAgent() {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token"})
	}

	userIDFromToken := claims["user_id"].(uint)
	sessionIDFromToken := claims["session_id"].(string)

	userID, err := sv.NewID(userIDFromToken)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	err = h.ValidateUsecase.Execute(sessionIDFromToken, userID)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid session"})
	}

	resp := fmt.Sprintf("hola %s, usted esta autorizado", name)
	return c.JSON(http.StatusOK, map[string]string{"message": resp})
}

func (h *UserEchoHandler) CreateUser(c echo.Context) error {
	var req d.UserCreateRequest
	// bindea el body del request (json) al dto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	fullname, err := v.NewFullName(req.Name, req.Lastname)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	username, err := v.NewUserName(req.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	dni, err := v.NewDni(req.Dni)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	email, err := v.NewEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	password, err := v.NewPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	// Map DTO to domain entity

	user := d.NewUser(
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
		return c.JSON(http.StatusConflict, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *UserEchoHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"message": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteUserUsecase.Execute(*id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "user deleted succesfully"})
}

func (h *UserEchoHandler) GetAllUsers(c echo.Context) error {
	users, err := h.GetAllUsersUsecase.Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var usersDTO []d.UserDTO
	for _, user := range users {
		userDTO := &d.UserDTO{
			ID:       user.ID().Id,
			Name:     user.FullName().Name,
			Lastname: user.FullName().Lastname,
			Username: user.Username().Username,
			Dni:      user.Dni().Dni,
			Email:    user.Email().Email,
			Password: user.Password().Password,
		}
		usersDTO = append(usersDTO, *userDTO)
	}

	return c.JSON(http.StatusOK, usersDTO)
}

func (h *UserEchoHandler) GetUserByID(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"message": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	user, err := h.GetUserByIDUsecase.Execute(*id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	userDTO := &d.UserDTO{
		ID:       user.ID().Id,
		Name:     user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Username: user.Username().Username,
		Dni:      user.Dni().Dni,
		Email:    user.Email().Email,
		Password: user.Password().Password,
	}
	return c.JSON(http.StatusOK, userDTO)
}
