package infrastructure

import (
	"github.com/labstack/echo/v4"
	userInfr "suffgo/internal/users/infrastructure"
)

func InitializeVoteEchoRouter(e *echo.Echo, handler *VoteEchoHandler) {
	voteGroup := e.Group("/v1/votes")

	voteGroup.Use(userInfr.AuthMiddleware)
	voteGroup.POST("", handler.CreateVote)
	voteGroup.DELETE("/:id", handler.DeleteVote)
	voteGroup.GET("", handler.GetAllVotes)
	voteGroup.GET("/:id", handler.GetVoteByID)
}
