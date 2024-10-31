package infrastructure

import (

	"github.com/labstack/echo/v4"
)

func InitializeUserEchoRouter(e *echo.Echo, handler *UserEchoHandler) {
	userGroup := e.Group("/v1/users")

	userGroup.POST("", handler.CreateUser)
	userGroup.DELETE("/:id", handler.DeleteUser)
	userGroup.GET("", handler.GetAllUsers)
	userGroup.GET("/:id", handler.GetUserByID)

	userGroup.POST("/login", handler.Login)

}
