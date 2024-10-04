package handlers

import (
	"net/http"
	"suffgo/internal/user/models"
	"suffgo/internal/user/usecases"

	"github.com/labstack/echo/v4"

	"github.com/labstack/gommon/log"
)

type userHttpHandler struct {
	userUsecase usecases.UserUsecase
}

func NewuserHttpHandler(userUsecase usecases.UserUsecase) UserHandler {
	return &userHttpHandler{
		userUsecase: userUsecase,
	}
}

// Create
func (h *userHttpHandler) RegisterUser(c echo.Context) error {
	reqBody := new(models.AddUserData)

	if err := c.Bind(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		return response(c, http.StatusBadRequest, "Bad request")
	}

	if err := h.userUsecase.UserDataRegister(reqBody); err != nil {
		return response(c, http.StatusInternalServerError, "Processing data failed")
	}

	return response(c, http.StatusOK, "Registrado correctamente")
}

// Retrieve
func (h *userHttpHandler) GetUserByID(c echo.Context) error {
	userID := c.Param("id")

	userData, err := h.userUsecase.GetUserByID(userID)
	if err != nil {
		return response(c, http.StatusInternalServerError, "User not found")
	}

	return c.JSON(http.StatusOK, userData)
}

// Retrieve all
func (h *userHttpHandler) GetAll(c echo.Context) error {
	users, err := h.userUsecase.GetAll()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// Update

// Delete
func (h *userHttpHandler) DeleteUser(c echo.Context) error {
	userID := c.Param("id")

	err := h.userUsecase.DeleteUser(userID)
	if err != nil {
		return response(c, http.StatusInternalServerError, "User not found")
	}

	return response(c, http.StatusOK, "Registro eliminado con exito")
}
