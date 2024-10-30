package infrastructure

import (
	"net/http"
	sv "suffgo/internal/shared/domain/valueObjects"
	v "suffgo/internal/vote/application/useCases"
	d "suffgo/internal/vote/domain"

	"github.com/labstack/echo/v4"
)

type VoteEchoHandler struct {
	CreateVoteUsecase *v.CreateUsecase
}

func NewVoterEchoHandler(
	createUC *v.CreateUsecase,
) *VoteEchoHandler {
	return &VoteEchoHandler{
		CreateVoteUsecase: createUC,
	}
}

func (h *VoteEchoHandler) CreateVote(c echo.Context) error {
	var req d.VoteCreateRequest
	// bindea el body del request (json) al dto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID, err := sv.NewID(req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	optionID, err := sv.NewID(req.OptionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	vote := d.NeweVote(
		nil,
		userID,
		optionID,
	)

	err = h.CreateVoteUsecase.Execute(*vote)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, req)
}
