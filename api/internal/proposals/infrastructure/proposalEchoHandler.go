package infrastructure

import (
	"net/http"
	"strconv"
	u "suffgo/internal/proposals/application/useCases"

	"github.com/labstack/echo/v4"

	d "suffgo/internal/proposals/domain"
	v "suffgo/internal/proposals/domain/valueObjects"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type ProposalEchoHandler struct {
	CreateProposalUsecase  *u.CreateUsecase
	GetAllProposalUseCase  *u.GetAllUsecase
	GetByIDProposalUseCase *u.GetByIDUsecase
	DeleteProposalUseCase  *u.DeleteUseCase
}

func NewProposalEchoHandler(
	createUC *u.CreateUsecase,
	getAllUC *u.GetAllUsecase,
	getByID *u.GetByIDUsecase,
	deleteUC *u.DeleteUseCase,
) *ProposalEchoHandler {
	return &ProposalEchoHandler{
		CreateProposalUsecase:  createUC,
		GetAllProposalUseCase:  getAllUC,
		GetByIDProposalUseCase: getByID,
		DeleteProposalUseCase:  deleteUC,
	}
}

func (h *ProposalEchoHandler) CreateProposal(c echo.Context) error {
	var req d.ProposalCreateRequest

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
		archive,
		*title,
		description,
	)

	err = h.CreateProposalUsecase.Execute(*proposal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, req)

}

func (h *ProposalEchoHandler) GetAllProposal(c echo.Context) error {
	proposal, err := h.GetAllProposalUseCase.Execute()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var proposalDTO []d.PorposalDTO

	for _, prop := range proposal {

		propDTO := &d.PorposalDTO{
			ID:          prop.ID().Id,
			Archive:     &prop.Archive().Archive,
			Title:       prop.Title().Title,
			Description: &prop.Description().Description,
		}
		proposalDTO = append(proposalDTO, *propDTO)
	}

	return c.JSON(http.StatusOK, proposalDTO)

}

func (h *ProposalEchoHandler) GetProposalByID(c echo.Context) error {

	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	proposal, err := h.GetByIDProposalUseCase.Execute(*id)

	if err != nil {
		if err.Error() == "proposal not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	proposalDTO := &d.PorposalDTO{
		ID:          proposal.ID().Id,
		Archive:     &proposal.Archive().Archive,
		Title:       proposal.Title().Title,
		Description: &proposal.Description().Description,
	}
	return c.JSON(http.StatusOK, proposalDTO)
}

func (h *ProposalEchoHandler) DeleteProposal(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteProposalUseCase.Execute(*id)

	if err != nil {
		if err.Error() == "proposal not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "Proposal deleted succesfully"})
}
