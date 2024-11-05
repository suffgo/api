package infrastructure

import (
	"github.com/labstack/echo/v4"
)

func InitializeUserEchoRouter(e *echo.Echo, handler *RoomEchoHandler) {
	roomGroup := e.Group("/v1/rooms")

	roomGroup.POST("", handler.CreateRoom)
	roomGroup.DELETE("/:id", handler.DeleteRoom)
	roomGroup.GET("", handler.GetAllRooms)
	roomGroup.GET("/:id", handler.GetRoomByID)
	roomGroup.GET("/myRooms/:id", handler.GetRoomsByAdmin)

}
