package infrastructure

import (
	u "suffgo/internal/rooms/application/useCases"
	userInfr "suffgo/internal/users/infrastructure"

	"github.com/labstack/echo/v4"
)

func InitializeRoomEchoRouter(e *echo.Echo, handler *RoomEchoHandler, getRoomByIDUsecase *u.GetByIDUsecase) {

	roomGroup := e.Group("/v1/rooms")
	roomGroup.Use(userInfr.AuthMiddleware)

	roomGroup.POST("", handler.CreateRoom)
	roomGroup.GET("", handler.GetAllRooms)
	roomGroup.GET("/:id", handler.GetRoomByID)
	roomGroup.GET("/myRooms", handler.GetRoomsByAdmin)
	roomGroup.POST("/restore/:id", handler.Restore)
	roomGroup.POST("/join", handler.JoinRoom)
	roomGroup.POST("/addUser", handler.AddSingleUser)
	roomGroup.GET("/ws/:room_id", handler.WsHandler)

	roomGroup.Use(AdminRoomMiddleware(getRoomByIDUsecase))
	roomGroup.PUT("/:id", handler.Update)
	roomGroup.DELETE("/:id", handler.DeleteRoom)
}
