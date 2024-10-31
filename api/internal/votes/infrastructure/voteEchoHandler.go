package infrastructure

import (
	"net/http"
	"strconv"
	sv "suffgo/internal/shared/domain/valueObjects"
	v "suffgo/internal/votes/application/useCases"
	d "suffgo/internal/votes/domain"

	se "suffgo/internal/shared/domain/errors"

	"github.com/labstack/echo/v4"
)

type VoteEchoHandler struct {
	CreateVoteUsecase  *v.CreateUsecase
	DeleteVoteUsecase  *v.DeleteUsecase
	GetAllVoteUsecase  *v.GetAllUsecase
	GetVoteByIDUsecase *v.GetByIDUsecase
}

func NewVoteEchoHandler(
	createUC *v.CreateUsecase,
	deleteUC *v.DeleteUsecase,
	getAllUC *v.GetAllUsecase,
	getByIDUC *v.GetByIDUsecase,
) *VoteEchoHandler {
	return &VoteEchoHandler{
		CreateVoteUsecase:  createUC,
		DeleteVoteUsecase:  deleteUC,
		GetAllVoteUsecase:  getAllUC,
		GetVoteByIDUsecase: getByIDUC,
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

func (h *VoteEchoHandler) DeleteVote(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteVoteUsecase.Execute(*id)

	if err != nil {
		if err.Error() == "vote not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "User deleted succesfully"})
}

func (h *VoteEchoHandler) GetAllVotes(c echo.Context) error {
	votes, err := h.GetAllVoteUsecase.Execute()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var votesDTO []d.VoteDTO
	for _, vote := range votes {
		voteDTO := &d.VoteDTO{
			ID:       vote.ID().Id,
			OptionID: vote.OptionID().Id,
			UserID:   vote.UserID().Id,
		}
		votesDTO = append(votesDTO, *voteDTO)
	}

	return c.JSON(http.StatusOK, votesDTO)
}

func (h *VoteEchoHandler) GetVoteByID(c echo.Context) error {

	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	vote, err := h.GetVoteByIDUsecase.Execute(*id)

	if err != nil {
		if err.Error() == "vote not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	voteDTO := &d.VoteDTO{
		ID:       vote.ID().Id,
		OptionID: vote.OptionID().Id,
		UserID:   vote.UserID().Id,
	}
	return c.JSON(http.StatusOK, voteDTO)
}
