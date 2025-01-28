package infrastructure

import (
	"github.com/labstack/echo/v4"
)

func InitializeUserEchoRouter(e *echo.Echo, handler *UserEchoHandler) {
	userGroup := e.Group("/v1/users")

	userGroup.POST("", handler.CreateUser)
	userGroup.GET("", handler.GetAllUsers)
	userGroup.GET("/:id", handler.GetUserByID)
	userGroup.GET("/email", handler.GetUserByEmail)

	userGroup.POST("/login", handler.Login)
	userGroup.POST("/newPassword", handler.ChangePassword)
	userGroup.Use(AuthMiddleware)
	userGroup.POST("/logout", handler.Logout)
	userGroup.DELETE("/:id", handler.DeleteUser)
	userGroup.POST("/restore/:id", handler.Restore)

	userGroup.Use(AuthMiddleware)
	userGroup.POST("/logout", handler.Logout)
	userGroup.GET("/auth", handler.CheckAuth) //200ok si esta autenticado
}
