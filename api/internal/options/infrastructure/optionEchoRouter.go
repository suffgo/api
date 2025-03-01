package infrastructure

import (
	"github.com/labstack/echo/v4"
	userInfr "suffgo/internal/users/infrastructure"
)

func InitializeOptionEchoRouter(e *echo.Echo, handler *OptionEchoHandler) {
	optionGroup := e.Group("/v1/options")

	optionGroup.Use(userInfr.AuthMiddleware)
	optionGroup.POST("", handler.CreateOption)
	optionGroup.DELETE("/:id", handler.DeleteOption)
	optionGroup.GET("", handler.GetAllOptions)
	optionGroup.GET("/:id", handler.GetOptionByID)
}
