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
	userGroup.Use(AuthMiddleware)
	userGroup.POST("/logout", handler.Logout)
	userGroup.DELETE("/:id", handler.DeleteUser)
	userGroup.POST("/restore/:id", handler.Restore)
	userGroup.GET("/byRoom/:id", handler.GetUsersByRoom)
	userGroup.PUT("/newPassword", handler.ChangePassword)

	userGroup.POST("/logout", handler.Logout)
	userGroup.GET("/auth", handler.CheckAuth)
	userGroup.PUT("", handler.Update)
}
