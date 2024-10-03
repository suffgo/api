package handlers

import "github.com/labstack/echo/v4"

type UserHandler interface {
	RegisterUser(c echo.Context) error
	GetUserByID(c echo.Context) error
	DeleteUser(c echo.Context) error
	GetAll(c echo.Context) error
}
