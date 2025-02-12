package infrastructure

import (
	"net/http"
	"strconv"
	u "suffgo/internal/rooms/application/useCases"
	sv "suffgo/internal/shared/domain/valueObjects"

	"github.com/labstack/echo/v4"
)

func AdminRoomMiddleware(getRoomByIDUsecase *u.GetByIDUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idParam := c.Param("id")
			idInput, err := strconv.ParseInt(idParam, 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room ID"})
			}

			roomID, _ := sv.NewID(uint(idInput))

			currentUser, err := GetUserIDFromSession(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			currentRoom, err := getRoomByIDUsecase.Execute(*roomID)
			if err != nil {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "room not found"})
			}

			if currentRoom.AdminID() != *currentUser {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "you are not allowed to delete this room"})
			}

			return next(c)
		}
	}
}
