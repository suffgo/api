package infrastructure

import (
	"suffgo/cmd/config"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitializeUserEchoRouter(e *echo.Echo, handler *UserEchoHandler) {
	userGroup := e.Group("/v1/users")

	userGroup.POST("", handler.CreateUser)
	userGroup.DELETE("/:id", handler.DeleteUser)
	userGroup.GET("", handler.GetAllUsers)
	userGroup.GET("/:id", handler.GetUserByID)

	userGroup.POST("/login", handler.Login)

	secureGroup := e.Group("/secure")
	secureGroup.Use(echojwt.JWT([]byte(config.SecretKey)))
	secureGroup.GET("", handler.SecureHello)
}
