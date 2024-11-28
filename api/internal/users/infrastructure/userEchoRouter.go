package infrastructure

import (

	"github.com/labstack/echo/v4"
)

func InitializeUserEchoRouter(e *echo.Echo, handler *UserEchoHandler) {
	userGroup := e.Group("/v1/users")

	userGroup.POST("", handler.CreateUser)
	userGroup.GET("", handler.GetAllUsers)
	userGroup.GET("/:id", handler.GetUserByID)

	userGroup.POST("/login", handler.Login)

	userGroup.Use(AuthMiddleware)
	userGroup.POST("/logout", handler.Logout)
	userGroup.DELETE("/:id", handler.DeleteUser)
	userGroup.GET("/auth", handler.CheckAuth) //200ok si esta autenticado
}
