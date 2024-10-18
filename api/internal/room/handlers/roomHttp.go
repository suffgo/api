package handlers

import (
	"net/http"
	"suffgo/internal/room/models"
	usecases "suffgo/internal/room/usecases"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type roomHttpHandler struct {
	roomUsecase usecases.RoomUsecase
}

func NewRoomHttpHandler(roomUsecase usecases.RoomUsecase) RoomHandler {
	return &roomHttpHandler{
		roomUsecase: roomUsecase,
	}
}

func (r *roomHttpHandler) RegisterRoom(c echo.Context) error {
	reqBody := new(models.AddRoomData)

	if err := c.Bind(reqBody); err != nil {
		log.Errorf("Error binding request body: %v", err)
		return response(c, http.StatusBadRequest, "Bad request")
	}

	if err := r.roomUsecase.RoomDataRegister(reqBody); err != nil {
		return response(c, http.StatusInternalServerError, "Processing data failed")
	}

	return response(c, http.StatusOK, "Room created succesfully")
}

func (r *roomHttpHandler) GetRoomByID(c echo.Context) error {
	roomID := c.Param("id")

	userData, err := r.roomUsecase.GetRoomByID(roomID)
	if err != nil {
		return response(c, http.StatusInternalServerError, "Room not found")
	}

	return c.JSON(http.StatusOK, userData)
}

func (r *roomHttpHandler) DeleteRoom(c echo.Context) error {
	return nil
}

func (r *roomHttpHandler) GetAll(c echo.Context) error {
	return nil
}
