package infrastructure

import (
	userInfr "suffgo/internal/users/infrastructure"

	"github.com/labstack/echo/v4"
)

func InitializeRoomEchoRouter(e *echo.Echo, handler *RoomEchoHandler) {

	roomGroup := e.Group("/v1/rooms")
	roomGroup.GET("", handler.GetAllRooms)

	roomGroup.Use(userInfr.AuthMiddleware)
	roomGroup.POST("", handler.CreateRoom)
	roomGroup.GET("/:id", handler.GetRoomByID)
	roomGroup.DELETE("/:id", handler.DeleteRoom)
	roomGroup.GET("/myRooms", handler.GetRoomsByAdmin)
	roomGroup.GET("/:id", handler.GetRoomByID)
	roomGroup.POST("/restore/:id", handler.Restore)
	roomGroup.POST("/join", handler.JoinRoom)
	roomGroup.POST("/addUser", handler.AddSingleUser)
	roomGroup.GET("/ws/:room_id", handler.WsHandler)
	roomGroup.PUT("/:id", handler.Update)
	roomGroup.DELETE("/whitelist/removeUser", handler.RemoveFromWhitelistHandler)
	roomGroup.GET("/history", handler.History)
}
