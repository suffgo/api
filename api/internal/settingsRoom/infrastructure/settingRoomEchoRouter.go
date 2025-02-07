package infrastructure

import (
	"github.com/labstack/echo/v4"
)

func InitializeSettingRoomEchoRouter(e *echo.Echo, handler *SettingRoomEchoHandler) {
	settingRoomGroup := e.Group("/v1/settingsRoom")

	settingRoomGroup.POST("", handler.CreateSettingRoom)
	settingRoomGroup.GET("", handler.GetAllSettingRoom)
	settingRoomGroup.GET("/:id", handler.GetSettingRoomByID)
	settingRoomGroup.DELETE("/:id", handler.DeleteSettingRoom)
	settingRoomGroup.PUT("/:id", handler.Update)

}
