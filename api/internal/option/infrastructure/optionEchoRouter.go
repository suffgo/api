package infrastructure

import (
	"github.com/labstack/echo/v4"
)

func InitializeOptionEchoRouter(e *echo.Echo, handler *OptionEchoHandler) {
	optionGroup := e.Group("/v1/options")

	optionGroup.POST("", handler.CreateOption)
	optionGroup.DELETE("/:id", handler.DeleteOption)
	optionGroup.GET("", handler.GetAllOptions)
	optionGroup.GET("/:id", handler.GetOptionByID)
	optionGroup.GET("/value/:value", handler.GetOptionByValue)
}
