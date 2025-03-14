package infrastructure

import (
	userInfr "suffgo/internal/users/infrastructure"

	"github.com/labstack/echo/v4"
)

func InitializeProposalEchoRouter(e *echo.Echo, handler *ProposalEchoHandler) {
	proposalGroup := e.Group("/v1/proposals")
	proposalGroup.GET("", handler.GetAllProposal)
	proposalGroup.GET("/:id", handler.GetProposalByID)

	proposalGroup.POST("/restore/:id", handler.RestoreProposal)	

	proposalGroup.Use(userInfr.AuthMiddleware)
	proposalGroup.GET("/byRoom/:room_id", handler.GetProposalsByRoomId)
	proposalGroup.POST("", handler.CreateProposal)
	proposalGroup.DELETE("/:id", handler.DeleteProposal)
	proposalGroup.PUT("/:id", handler.Update)
}
