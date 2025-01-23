package infrastructure

import (
	"errors"
	"net/http"
	"strconv"
	r "suffgo/internal/rooms/application/useCases"

	d "suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"

	sv "suffgo/internal/shared/domain/valueObjects"

	se "suffgo/internal/shared/domain/errors"

	rerr "suffgo/internal/rooms/domain/errors"

	"github.com/labstack/echo/v4"
)

type RoomEchoHandler struct {
	CreateRoomUsecase  *r.CreateUsecase
	DeleteRoomUsecase  *r.DeleteUsecase
	GetAllUsecase      *r.GetAllUsecase
	GetRoomByIDUsecase *r.GetByIDUsecase
	GetByAdminUsecase  *r.GetByAdminUsecase
	RestoreUsecase     *r.RestoreUsecase
	JoinRoomUsecase    *r.JoinRoomUsecase
}

func NewRoomEchoHandler(
	creatUC *r.CreateUsecase,
	deleteUC *r.DeleteUsecase,
	getAllUC *r.GetAllUsecase,
	getByIDUC *r.GetByIDUsecase,
	getByAdminUC *r.GetByAdminUsecase,
	restoreUC *r.RestoreUsecase,
	joinRoomUC *r.JoinRoomUsecase,
) *RoomEchoHandler {
	return &RoomEchoHandler{
		CreateRoomUsecase:  creatUC,
		DeleteRoomUsecase:  deleteUC,
		GetAllUsecase:      getAllUC,
		GetRoomByIDUsecase: getByIDUC,
		GetByAdminUsecase:  getByAdminUC,
		RestoreUsecase:     restoreUC,
		JoinRoomUsecase:    joinRoomUC,
	}
}

func (h *RoomEchoHandler) CreateRoom(c echo.Context) error {
	var req d.RoomCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	linkInvite, err := v.NewLinkInvite(req.LinkInvite)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	isFormal, err := v.NewIsFormal(req.IsFormal)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	name, err := v.NewName(req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Obtener el user_id de la sesion
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuario no autenticado"})
	}

	adminIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	adminID, err := sv.NewID(uint(adminIDUint))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	room := d.NewRoom(
		nil,
		*linkInvite,
		*isFormal,
		*name,
		adminID,
	)

	createdRoom, err := h.CreateRoomUsecase.Execute(*room)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	roomDTO := &d.RoomDTO{
		ID:         createdRoom.ID().Id,
		LinkInvite: createdRoom.LinkInvite().LinkInvite,
		IsFormal:   createdRoom.IsFormal().IsFormal,
		Name:       createdRoom.Name().Name,
		AdminID:    createdRoom.AdminID().Id,
		RoomCode:   createdRoom.InviteCode().Code,
	}

	response := map[string]interface{}{
		"success": "éxito al crear sala",
		"room":    roomDTO,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *RoomEchoHandler) DeleteRoom(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteRoomUsecase.Execute(*id)
	if err != nil {
		if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"success": "room deleted succesfully"})
}

func (h *RoomEchoHandler) GetAllRooms(c echo.Context) error {
	rooms, err := h.GetAllUsecase.Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var roomsDTO []d.RoomDTO
	for _, room := range rooms {
		roomDTO := &d.RoomDTO{
			ID:         room.ID().Id,
			LinkInvite: room.LinkInvite().LinkInvite,
			IsFormal:   room.IsFormal().IsFormal,
			Name:       room.Name().Name,
			AdminID:    room.AdminID().Id,
		}
		roomsDTO = append(roomsDTO, *roomDTO)
	}
	return c.JSON(http.StatusOK, roomsDTO)
}

func (h *RoomEchoHandler) GetRoomByID(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrDataMap.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	room, err := h.GetRoomByIDUsecase.Execute(*id)
	if err != nil {
		if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	roomDTO := &d.RoomDTO{
		ID:         room.ID().Id,
		LinkInvite: room.LinkInvite().LinkInvite,
		IsFormal:   room.IsFormal().IsFormal,
		Name:       room.Name().Name,
		AdminID:    room.AdminID().Id,
	}
	return c.JSON(http.StatusOK, roomDTO)
}

func (h *RoomEchoHandler) GetRoomsByAdmin(c echo.Context) error {

	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuario no autenticado"})
	}

	idInt, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	userID, err := sv.NewID(uint(idInt))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	rooms, err := h.GetByAdminUsecase.Execute(*userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var roomsDTO []d.RoomDTO
	for _, room := range rooms {
		roomDTO := &d.RoomDTO{
			ID:         room.ID().Id,
			LinkInvite: room.LinkInvite().LinkInvite,
			IsFormal:   room.IsFormal().IsFormal,
			Name:       room.Name().Name,
			AdminID:    room.AdminID().Id,
		}
		roomsDTO = append(roomsDTO, *roomDTO)
	}
	return c.JSON(http.StatusOK, roomsDTO)
}

func (h *RoomEchoHandler) JoinRoom(c echo.Context) error {
	var req d.JoinRoomRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	room, err := h.JoinRoomUsecase.Execute(req.RoomCode)

	if err != nil {
		if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	roomDTO := &d.RoomDTO{
		ID:         room.ID().Id,
		LinkInvite: room.LinkInvite().LinkInvite,
		IsFormal:   room.IsFormal().IsFormal,
		Name:       room.Name().Name,
		AdminID:    room.AdminID().Id,
		RoomCode:   room.InviteCode().Code,
	}

	response := map[string]interface{}{
		"success": "Ingreso a la sala exitoso",
		"room":    roomDTO,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *RoomEchoHandler) Restore(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.RestoreUsecase.Execute(*id)
	if err != nil {
		if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Room not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"succes": "room restored succesfully"})

}
