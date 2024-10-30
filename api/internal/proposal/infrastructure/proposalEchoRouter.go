package infrastructure

import (
	"github.com/labstack/echo/v4"
)

func InitializeProposalEchoRouter(e *echo.Echo, handler *ProposalEchoHandler) {
	proposalGroup := e.Group("/v1/proposals")

	proposalGroup.POST("", handler.CreateProposal)
	proposalGroup.DELETE("/:id", handler.DeleteProposal)
	proposalGroup.GET("", handler.GetAllProposal)
	proposalGroup.GET("/:id", handler.GetProposalByID)
}
