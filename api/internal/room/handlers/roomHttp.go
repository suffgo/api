package handlers

import (
	usecases "suffgo/internal/room/usescases"

	"github.com/labstack/echo/v4"
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
	return nil
}

func (r *roomHttpHandler) GetRoomByID(c echo.Context) error {
	return nil
}

func (r *roomHttpHandler) DeleteRoom(c echo.Context) error {
	return nil
}

func (r *roomHttpHandler) GetAll(c echo.Context) error {
	return nil
}
