package infrastructure

import (
	userInfr "suffgo/internal/users/infrastructure"

	"github.com/labstack/echo/v4"
)

func InitializeRoomEchoRouter(e *echo.Echo, handler *RoomEchoHandler) {

	roomGroup := e.Group("/v1/rooms")
	roomGroup.Use(userInfr.AuthMiddleware)

	roomGroup.POST("", handler.CreateRoom)
	roomGroup.DELETE("/:id", handler.DeleteRoom)
	roomGroup.GET("", handler.GetAllRooms)
	roomGroup.GET("/:id", handler.GetRoomByID)
	roomGroup.GET("/myRooms", handler.GetRoomsByAdmin)
	roomGroup.POST("/restore/:id", handler.Restore)
	roomGroup.POST("/join", handler.JoinRoom)
	roomGroup.POST("/addUser", handler.AddSingleUser)
	roomGroup.GET("/ws/:room_id", handler.WsHandler)
	roomGroup.PUT("/:id", handler.Update)
}
