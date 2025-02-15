package infrastructure

import (
	"errors"
	"net/http"
	"strconv"
	u "suffgo/internal/proposals/application/useCases"

	"github.com/labstack/echo/v4"

	d "suffgo/internal/proposals/domain"
	v "suffgo/internal/proposals/domain/valueObjects"
	rh "suffgo/internal/rooms/infrastructure"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"

	perrors "suffgo/internal/proposals/domain/errors"
)

type ProposalEchoHandler struct {
	CreateProposalUsecase  *u.CreateUsecase
	GetAllProposalUseCase  *u.GetAllUsecase
	GetByIDProposalUseCase *u.GetByIDUsecase
	DeleteProposalUseCase  *u.DeleteUseCase
	RestoreUseCase         *u.RestoreUsecase
	UpdateUseCase          *u.UpdateUsecase
}

func NewProposalEchoHandler(
	createUC *u.CreateUsecase,
	getAllUC *u.GetAllUsecase,
	getByID *u.GetByIDUsecase,
	deleteUC *u.DeleteUseCase,
	restoreUC *u.RestoreUsecase,
	updateUC *u.UpdateUsecase,
) *ProposalEchoHandler {
	return &ProposalEchoHandler{
		CreateProposalUsecase:  createUC,
		GetAllProposalUseCase:  getAllUC,
		GetByIDProposalUseCase: getByID,
		DeleteProposalUseCase:  deleteUC,
		RestoreUseCase:         restoreUC,
		UpdateUseCase:          updateUC,
	}
}

func (h *ProposalEchoHandler) CreateProposal(c echo.Context) error {
	var req d.ProposalCreateRequest

	//obtengo id del usuario de la sesion para verificar que sea el dueño de la sala
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuario no autenticado"})
	}

	userCreatorIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id de usuario inválido"})
	}

	userCreatorID, err := sv.NewID(uint(userCreatorIDUint))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

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

	roomID, err := sv.NewID(req.RoomID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	proposal := d.NewProposal(
		nil,
		archive,
		*title,
		description,
		roomID,
	)

	createdProp, err := h.CreateProposalUsecase.Execute(*proposal, *userCreatorID)
	if err != nil {

		if err.Error() == "operación no autorizada para este usuario" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		} else if errors.Is(err, se.ErrInvalidID) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	proposalDTO := d.ProposalDTO{
		ID:          createdProp.ID().Id,
		Archive:     &createdProp.Archive().Archive,
		Title:       createdProp.Title().Title,
		Description: &createdProp.Description().Description,
		RoomID:      createdProp.RoomID().Id,
	}

	response := map[string]interface{}{
		"success":  "éxito al crear propuesta",
		"proposal": proposalDTO,
	}

	return c.JSON(http.StatusCreated, response)

}

func (h *ProposalEchoHandler) GetAllProposal(c echo.Context) error {
	proposal, err := h.GetAllProposalUseCase.Execute()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var proposalDTO []d.ProposalDTO

	for _, prop := range proposal {

		propDTO := &d.ProposalDTO{
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	proposal, err := h.GetByIDProposalUseCase.Execute(*id)

	if err != nil {
		if errors.Is(err, perrors.ErrPropNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	proposalDTO := &d.ProposalDTO{
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))

	currentUser, err := rh.GetUserIDFromSession(c)
	if err != nil {
		return err
	}

	err = h.DeleteProposalUseCase.Execute(*id, *currentUser)

	if err != nil {
		if errors.Is(err, perrors.ErrPropNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "Proposal deleted succesfully"})
}

func (h *ProposalEchoHandler) RestoreProposal(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.RestoreUseCase.Execute(*id)

	if err != nil {
		if errors.Is(err, perrors.ErrPropNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Proposal not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"succes": "proposal restored succesfully"})

}

func (h *ProposalEchoHandler) Update(c echo.Context) error {
	proposalIDStr := c.Param("id")
	proposalID, err := strconv.Atoi(proposalIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid proposal ID"})
	}

	var req d.ProposalUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	id, err := sv.NewID(uint(proposalID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	currentProposal, err := h.GetByIDProposalUseCase.Execute(*id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	Archive, err := v.NewArchive(*req.Archive)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	Title, err := v.NewTitle(req.Title)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	Description, err := v.NewDescription(*req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	RoomID := currentProposal.RoomID()

	proposal := d.NewProposal(
		id,
		Archive,
		*Title,
		Description,
		&RoomID,
	)

	currentUser, err := rh.GetUserIDFromSession(c)
	if err != nil {
		return err
	}

	updatedProposal, err := h.UpdateUseCase.Execute(proposal, *currentUser)

	if err != nil {
		if errors.Is(err, perrors.ErrPropNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	proposalDTO := d.ProposalDTO{
		ID:          updatedProposal.ID().Id,
		Archive:     &updatedProposal.Archive().Archive,
		Title:       updatedProposal.Title().Title,
		Description: &updatedProposal.Description().Description,
		RoomID:      updatedProposal.RoomID().Id,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":     "proposal updated successfully",
		"settingRoom": proposalDTO,
	})
}
