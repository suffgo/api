package handlers

import "github.com/labstack/echo/v4"

type RoomHandler interface {
	RegisterRoom(c echo.Context) error
	GetRoomByID(c echo.Context) error
	DeleteRoom(c echo.Context) error
	GetAll(c echo.Context) error
}
