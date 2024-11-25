package infrastructure

import (
	"github.com/labstack/echo/v4"
	userInfr "suffgo/internal/users/infrastructure"
)

func InitializeProposalEchoRouter(e *echo.Echo, handler *ProposalEchoHandler) {
	proposalGroup := e.Group("/v1/proposals")
	proposalGroup.GET("", handler.GetAllProposal)
	proposalGroup.GET("/:id", handler.GetProposalByID)

	proposalGroup.Use(userInfr.AuthMiddleware)
	proposalGroup.POST("", handler.CreateProposal)
	proposalGroup.DELETE("/:id", handler.DeleteProposal)
}