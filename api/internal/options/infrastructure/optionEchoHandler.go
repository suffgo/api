package infrastructure

import (
	"errors"
	"net/http"
	"strconv"
	u "suffgo/internal/options/application/useCases"

	d "suffgo/internal/options/domain"
	v "suffgo/internal/options/domain/valueObjects"

	"github.com/labstack/echo/v4"

	sv "suffgo/internal/shared/domain/valueObjects"

	oerrors "suffgo/internal/options/domain/errors"
	se "suffgo/internal/shared/domain/errors"
)

type OptionEchoHandler struct {
	CreateOptionUsecase     *u.CreateUsecase
	DeleteOptionUsecase     *u.DeleteUsecase
	GetAllOptionUsecase     *u.GetAllUsecase
	GetOtionByIDUsecase     *u.GetByIDUsecase
	GetOptionByValueUsecase *u.GetByValueUsecase
	GetOptionByPropUsecase  *u.GetByPropUsecase
}

func NewOptionEchoHandler(
	createUC *u.CreateUsecase,
	deleteUC *u.DeleteUsecase,
	getAllUC *u.GetAllUsecase,
	getByIDUC *u.GetByIDUsecase,
	getByValue *u.GetByValueUsecase,
	getByProp *u.GetByPropUsecase,
) *OptionEchoHandler {
	return &OptionEchoHandler{
		CreateOptionUsecase:     createUC,
		DeleteOptionUsecase:     deleteUC,
		GetAllOptionUsecase:     getAllUC,
		GetOtionByIDUsecase:     getByIDUC,
		GetOptionByValueUsecase: getByValue,
		GetOptionByPropUsecase:  getByProp,
	}
}

func (h *OptionEchoHandler) CreateOption(c echo.Context) error {
	var req d.OptionCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	value, err := v.NewValue(req.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	proposalID, err := sv.NewID(req.ProposalID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	option := d.NewOption(
		nil,
		*value,
		proposalID,
	)

	err = h.CreateOptionUsecase.Execute(*option)
	if err != nil {
		if errors.Is(err, oerrors.ErrOptRepeated) {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *OptionEchoHandler) DeleteOption(c echo.Context) error {

	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteOptionUsecase.Execute(*id)

	if err != nil {
		if errors.Is(err, oerrors.ErrOptNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"succes": "Option deleted succesfully"})
}

func (h *OptionEchoHandler) GetAllOptions(c echo.Context) error {
	options, err := h.GetAllOptionUsecase.Execute()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var OptionsDTO []d.OptionDTO
	for _, option := range options {
		OptionDTO := &d.OptionDTO{
			ID:         option.ID().Id,
			Value:      option.Value().Value,
			ProposalID: option.ProposalID().Id,
		}
		OptionsDTO = append(OptionsDTO, *OptionDTO)
	}

	return c.JSON(http.StatusOK, OptionsDTO)
}

func (h *OptionEchoHandler) GetOptionByID(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	option, err := h.GetOtionByIDUsecase.Execute(*id)

	if err != nil {
		if errors.Is(err, oerrors.ErrOptNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	optionDTO := &d.OptionDTO{
		ID:         option.ID().Id,
		Value:      option.Value().Value,
		ProposalID: option.ProposalID().Id,
	}

	return c.JSON(http.StatusOK, optionDTO)
}

func (h *OptionEchoHandler) GetOptionByValue(c echo.Context) error {
	valueParam := c.Param("value")

	value, err := v.NewValue(valueParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}
	option, err := h.GetOptionByValueUsecase.Execute(*value)

	if err != nil {
		if errors.Is(err, oerrors.ErrOptNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	optionDTO := &d.OptionDTO{
		ID:         option.ID().Id,
		Value:      option.Value().Value,
		ProposalID: option.ProposalID().Id,
	}
	return c.JSON(http.StatusOK, optionDTO)

}

func (h *OptionEchoHandler) GetOptionByProposal(c echo.Context) error {

	idParam := c.Param("proposal_id")

	propId, err := sv.NewID(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	options, err := h.GetOptionByPropUsecase.Execute(propId)

	if err != nil {
		if errors.Is(oerrors.ErrOptNotFound, err) {
			return c.JSON(http.StatusNoContent, map[string]string{"error": se.ErrInvalidID.Error()})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	var optionsDto []d.OptionDTO
	for _, option := range options {
		OptionDTO := &d.OptionDTO{
			ID:         option.ID().Id,
			Value:      option.Value().Value,
			ProposalID: option.ProposalID().Id,
		}
		optionsDto = append(optionsDto, *OptionDTO)
	}

	response := map[string]interface{}{
		"success": "opciones de la propuesta " + idParam,
		"options": optionsDto,
	}

	return c.JSON(http.StatusOK, response)
}
