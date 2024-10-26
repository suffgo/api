package infrastructure

import (
	"net/http"
	u "suffgo/internal/proposal/application/useCases"

	"github.com/labstack/echo/v4"

	d "suffgo/internal/proposal/domain"
	v "suffgo/internal/proposal/domain/valueObjects"
)

type ProposalEchoHandler struct {
	CreateProposalUsecase *u.CreateUsecase
}

func NewProposalEchoHandler(createProposalUC *u.CreateUsecase) *ProposalEchoHandler {
	return &ProposalEchoHandler{
		CreateProposalUsecase: createProposalUC,
	}
}

func (h *ProposalEchoHandler) CreateProposal(c echo.Context) error {
	var req d.PorposalCreateRequest

	// bindea el body del request (json) al dto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	archive, err := v.NewArchive(*req.Archive)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	title, err := v.NewTitle(req.Title)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	description, err := v.NewDescription(*req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	proposal := d.NewProposal(
		nil,
		*archive,
		*title,
		*description,
	)

	err = h.CreateProposalUsecase.Execute(*proposal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, req)

}
