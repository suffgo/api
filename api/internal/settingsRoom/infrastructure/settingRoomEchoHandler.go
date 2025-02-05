package infrastructure

import (
	"errors"
	"net/http"
	"strconv"
	s "suffgo/internal/settingsRoom/application/useCases"
	d "suffgo/internal/settingsRoom/domain"
	seterr "suffgo/internal/settingsRoom/domain/errors"
	v "suffgo/internal/settingsRoom/domain/valueObjects"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"

	"github.com/labstack/echo/v4"
)

type SettingRoomEchoHandler struct {
	CreateSettingRoomUsecase  *s.CreateUsecase
	DeleteSettingRoomUsecase  *s.DeleteUsecase
	GetAllSettingRoomUsecase  *s.GetAllUsecase
	GetSettingRoomByIDUsecase *s.GetByIDUsecase
}

func NewSettingRoomEchoHandler(
	createUC *s.CreateUsecase,
	deleteUC *s.DeleteUsecase,
	getAllUC *s.GetAllUsecase,
	getByIDUC *s.GetByIDUsecase,
) *SettingRoomEchoHandler {
	return &SettingRoomEchoHandler{
		CreateSettingRoomUsecase:  createUC,
		DeleteSettingRoomUsecase:  deleteUC,
		GetAllSettingRoomUsecase:  getAllUC,
		GetSettingRoomByIDUsecase: getByIDUC,
	}
}

func (h *SettingRoomEchoHandler) CreateSettingRoom(c echo.Context) error {
	var req d.SettingRoomCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	privacy, err := v.NewPrivacy(req.Privacy)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	proposalTimer, err := v.NewProposalTimer(req.ProposalTimer)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	quorum, err := v.NewQuorum(req.Quorum)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	timeAndDate, err := v.NewDateTime(req.DateTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}	

	voterLimit, err := v.NewVoterLimit(req.VoterLimit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	roomID, err := sv.NewID(req.RoomID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// Generar un nuevo ID para el SettingRoom
	newID, err := sv.NewID(0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// asegurarnos de que newID no sea nil antes de crear el SettingRoom
	if newID == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "could not generate new ID"})
	}

	settingRoom := d.NewSettingRoom(
		newID,
		privacy,
		proposalTimer,
		*quorum,
		*timeAndDate,
		voterLimit,
		roomID,
	)

	if settingRoom == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "could not create setting room"})
	}

	err = h.CreateSettingRoomUsecase.Execute(*settingRoom)
	if err != nil {

		if errors.Is(err, seterr.ErrAlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *SettingRoomEchoHandler) DeleteSettingRoom(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteSettingRoomUsecase.Execute(*id)
	if err != nil {
		if err.Error() == "setting room not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "setting room deleted succesfully"})
}

func (h *SettingRoomEchoHandler) GetAllSettingRoom(c echo.Context) error {
	settingsRoom, err := h.GetAllSettingRoomUsecase.Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var settingsRoomDTO []d.SettingRoomDTO
	for _, settingRoom := range settingsRoom {
		SettingRoomDTO := &d.SettingRoomDTO{
			ID:            settingRoom.ID().Id,
			Privacy:       settingRoom.Privacy().Privacy,
			ProposalTimer: settingRoom.ProposalTimer().ProposalTimer,
			Quorum:        settingRoom.Quorum().Quorum,
			StartTime:     settingRoom.StartTime().DateTime,
			VoterLimit:    settingRoom.VoterLimit().VoterLimit,
			RoomID:        settingRoom.RoomID().Id,
		}
		settingsRoomDTO = append(settingsRoomDTO, *SettingRoomDTO)
	}

	return c.JSON(http.StatusOK, settingsRoomDTO)
}

func (h *SettingRoomEchoHandler) GetSettingRoomByID(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	settingRoom, err := h.GetSettingRoomByIDUsecase.Execute(*id)

	if err != nil {
		if err.Error() == "setting room not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	settingRoomDTO := &d.SettingRoomDTO{
		ID:            settingRoom.ID().Id,
		Privacy:       settingRoom.Privacy().Privacy,
		ProposalTimer: settingRoom.ProposalTimer().ProposalTimer,
		Quorum:        settingRoom.Quorum().Quorum,
		StartTime:     settingRoom.StartTime().DateTime,
		VoterLimit:    settingRoom.VoterLimit().VoterLimit,
		RoomID:        settingRoom.RoomID().Id,
	}
	return c.JSON(http.StatusOK, settingRoomDTO)
}
