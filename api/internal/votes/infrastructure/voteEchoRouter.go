package infrastructure

import (
	"github.com/labstack/echo/v4"
)

func InitializeVoteEchoRouter(e *echo.Echo, handler *VoteEchoHandler) {
	voteGroup := e.Group("/v1/vote")

	voteGroup.POST("", handler.CreateVote)
	voteGroup.DELETE("/:id", handler.DeleteVote)
	voteGroup.GET("", handler.GetAllVotes)
	voteGroup.GET("/:id", handler.GetVoteByID)
}
