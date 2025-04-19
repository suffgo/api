package infrastructure

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	r "suffgo/internal/rooms/application/useCases"
	addUsers "suffgo/internal/rooms/application/useCases/addUsers"
	roomWs "suffgo/internal/rooms/application/useCases/websocket"

	d "suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"

	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"

	rerr "suffgo/internal/rooms/domain/errors"
	uerr "suffgo/internal/users/domain/errors"

	useruc "suffgo/internal/users/application/useCases"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"suffgo/internal/settingsRoom/domain"
	srerr "suffgo/internal/settingsRoom/domain/errors"
)

type RoomEchoHandler struct {
	CreateRoomUsecase    *r.CreateUsecase
	DeleteRoomUsecase    *r.DeleteUsecase
	GetAllUsecase        *r.GetAllUsecase
	GetRoomByIDUsecase   *r.GetByIDUsecase
	GetByAdminUsecase    *r.GetByAdminUsecase
	RestoreUsecase       *r.RestoreUsecase
	JoinRoomUsecase      *r.JoinRoomUsecase
	AddSingleUserUsecase *addUsers.AddSingleUserUsecase
	GetUserByIDUsecase   *useruc.GetByIDUsecase
	UpdateRoomUsecase    *r.UpdateRoomUsecase
	GetSrByRoomIDUsecase *r.GetSrByRoomUsecase
	ManageWsUsecase      *roomWs.ManageWsUsecase
	WhiteListRmUsecase   *r.WhitelistRmUsecase
	HistoryRoomsUsecase  *r.HistoryRooms
}

func NewRoomEchoHandler(
	creatUC *r.CreateUsecase,
	deleteUC *r.DeleteUsecase,
	getAllUC *r.GetAllUsecase,
	getByIDUC *r.GetByIDUsecase,
	getByAdminUC *r.GetByAdminUsecase,
	restoreUC *r.RestoreUsecase,
	joinRoomUC *r.JoinRoomUsecase,
	addSingleUserUC *addUsers.AddSingleUserUsecase,
	getUserByIDUC *useruc.GetByIDUsecase,
	updateUC *r.UpdateRoomUsecase,
	manageWsUC *roomWs.ManageWsUsecase,
	getSrByRoomIDUC *r.GetSrByRoomUsecase,
	whitelistRmUC *r.WhitelistRmUsecase,
	historyRoomsUC *r.HistoryRooms,

) *RoomEchoHandler {
	return &RoomEchoHandler{
		CreateRoomUsecase:    creatUC,
		DeleteRoomUsecase:    deleteUC,
		GetAllUsecase:        getAllUC,
		GetRoomByIDUsecase:   getByIDUC,
		GetByAdminUsecase:    getByAdminUC,
		RestoreUsecase:       restoreUC,
		JoinRoomUsecase:      joinRoomUC,
		AddSingleUserUsecase: addSingleUserUC,
		GetUserByIDUsecase:   getUserByIDUC,
		UpdateRoomUsecase:    updateUC,
		ManageWsUsecase:      manageWsUC,
		GetSrByRoomIDUsecase: getSrByRoomIDUC,
		WhiteListRmUsecase:   whitelistRmUC,
		HistoryRoomsUsecase:  historyRoomsUC,
	}
}

func (h *RoomEchoHandler) CreateRoom(c echo.Context) error {
	var req d.RoomCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Datos de entrada inválidos: " + err.Error()})
	}

	isFormal, err := v.NewIsFormal(req.IsFormal)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tipo de sala inválido"})
	}

	name, err := v.NewName(req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre inválido"})
	}

	description, err := v.NewDescription(req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Descripción inválida"})
	}

	image, err := v.NewImage(req.Image)
	if err != nil {
		log.Printf("why %s\n", err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Imagen Invalida"})
	}

	// Obtener el ID del administrador desde la sesión
	adminID, err := GetUserIDFromSession(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Usuario no autenticado"})
	}
	// Crear objeto de sala
	state, _ := v.NewState("created")
	room := d.NewRoom(nil, *isFormal, nil, *name, adminID, *description, image, state)

	// Ejecutar caso de uso
	createdRoom, err := h.CreateRoomUsecase.Execute(*room)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Error al crear la sala: " + err.Error()})
	}

	// Crear respuesta con DTO
	roomDTO := &d.RoomDTO{
		ID:          createdRoom.ID().Id,
		IsFormal:    createdRoom.IsFormal().IsFormal,
		Name:        createdRoom.Name().Name,
		AdminID:     createdRoom.AdminID().Id,
		Description: createdRoom.Description().Description,
		Code:        createdRoom.Code().Code,
		State:       createdRoom.State().CurrentState,
		Image:       createdRoom.Image().URL(),
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": "Sala creada con éxito",
		"room":    roomDTO,
	})
}

func (h *RoomEchoHandler) DeleteRoom(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))

	userID, err := GetUserIDFromSession(c)
	if err != nil {
		return err
	}
	err = h.DeleteRoomUsecase.Execute(*id, *userID)
	if err != nil {
		if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": err.Error()})
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

	var roomsDTO []d.RoomDetailedDTO
	for _, room := range rooms {

		admin, err := h.GetUserByIDUsecase.Execute(room.AdminID())
		if err != nil {
			if !errors.Is(err, uerr.ErrUserNotFound) {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var adminName string
		if admin != nil {
			adminName = admin.FullName().Name + " " + admin.FullName().Lastname
		} else {
			adminName = "null"
		}

		var roomDTO *d.RoomDetailedDTO
		if room.IsFormal().IsFormal {
			settingRoom, err := h.GetSrByRoomIDUsecase.Execute(room.ID())
			if err != nil {
				if errors.Is(err, srerr.SettingRoomNotFoundError) {
					return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			roomDTO = &d.RoomDetailedDTO{
				ID:          room.ID().Id,
				IsFormal:    room.IsFormal().IsFormal,
				RoomTitle:   room.Name().Name,
				AdminName:   adminName,
				Description: room.Description().Description,
				Code:        room.Code().Code,
				StartTime:   settingRoom.StartTime().DateTime,
				State:       room.State().CurrentState,
				Image:       room.Image().URL(),
			}
		} else {
			roomDTO = &d.RoomDetailedDTO{
				ID:          room.ID().Id,
				RoomTitle:   room.Name().Name,
				AdminName:   adminName,
				Description: room.Description().Description,
				Code:        room.Code().Code,
				State:       room.State().CurrentState,
				Image:       room.Image().URL(),
			}
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

	id, err := sv.NewID(uint(idInput))
	if err != nil || id == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID inválido"})
	}
	room, err := h.GetRoomByIDUsecase.Execute(*id)

	if err != nil {
		if errors.Is(rerr.ErrRoomNotFound, err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	admin, err := h.GetUserByIDUsecase.Execute(room.AdminID())

	if err != nil {
		if !errors.Is(uerr.ErrUserNotFound, err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
	}

	userId, _ := GetUserIDFromSession(c)

	var adminName string
	if admin == nil { //el admin es un usuario eliminado
		adminName = "null"
	} else {
		adminName = admin.FullName().Lastname + " " + admin.FullName().Name
	}

	privileges := false
	if room.AdminID().Id == userId.Id {
		privileges = true
	}

	var settingRoom *domain.SettingRoom
	var roomDetailedDTO *d.RoomDetailedDTO
	if room.IsFormal().IsFormal {
		settingRoom, err = h.GetSrByRoomIDUsecase.Execute(room.ID())
		if err != nil {
			if errors.Is(err, srerr.SettingRoomNotFoundError) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		roomDetailedDTO = &d.RoomDetailedDTO{
			ID:          room.ID().Id,
			IsFormal:    room.IsFormal().IsFormal,
			RoomTitle:   room.Name().Name,
			AdminName:   adminName,
			Description: room.Description().Description,
			Code:        room.Code().Code,
			StartTime:   settingRoom.StartTime().DateTime,
			State:       room.State().CurrentState,
			Image:       room.Image().URL(),
			Privileges:  privileges,
		}
	} else {
		roomDetailedDTO = &d.RoomDetailedDTO{
			ID:          room.ID().Id,
			IsFormal:    room.IsFormal().IsFormal,
			RoomTitle:   room.Name().Name,
			AdminName:   adminName,
			Description: room.Description().Description,
			Code:        room.Code().Code,
			State:       room.State().CurrentState,
			Image:       room.Image().URL(),
			Privileges:  privileges,
		}
	}

	response := map[string]interface{}{
		"success": "Ingreso a la sala exitoso",
		"room":    roomDetailedDTO,
	}

	return c.JSON(http.StatusOK, response)
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

	admin, err := h.GetUserByIDUsecase.Execute(*userID)

	if err != nil {
		if !errors.Is(err, uerr.ErrUserNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	var adminName string
	if admin != nil {
		adminName = admin.FullName().Name + " " + admin.FullName().Lastname
	} else {
		adminName = "null"
	}

	var roomsDTO []d.RoomDetailedDTO
	for _, room := range rooms {

		userId, _ := GetUserIDFromSession(c)
		privileges := false
		if userId.Id == room.AdminID().Id {
			privileges = true
		}

		var roomDTO *d.RoomDetailedDTO
		if room.IsFormal().IsFormal {
			settingRoom, err := h.GetSrByRoomIDUsecase.Execute(room.ID())

			if err != nil {
				if errors.Is(err, srerr.SettingRoomNotFoundError) {
					return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
				}
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			roomDTO = &d.RoomDetailedDTO{
				ID:          room.ID().Id,
				IsFormal:    room.IsFormal().IsFormal,
				RoomTitle:   room.Name().Name,
				AdminName:   adminName,
				Description: room.Description().Description,
				Code:        room.Code().Code,
				StartTime:   settingRoom.StartTime().DateTime,
				State:       room.State().CurrentState,
				Image:       room.Image().URL(),
				Privileges:  privileges,
			}
		} else {
			roomDTO = &d.RoomDetailedDTO{
				ID:          room.ID().Id,
				RoomTitle:   room.Name().Name,
				AdminName:   adminName,
				Description: room.Description().Description,
				Code:        room.Code().Code,
				State:       room.State().CurrentState,
				Image:       room.Image().URL(),
				Privileges:  privileges,
			}
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

	userID, err := GetUserIDFromSession(c)

	if err != nil {
		return err
	}

	room, err := h.JoinRoomUsecase.Execute(req.RoomCode, *userID)

	if err != nil {
		if errors.Is(err, rerr.ErrNotWhitelist) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	admin, err := h.GetUserByIDUsecase.Execute(room.AdminID())

	if err != nil {
		if !errors.Is(err, uerr.ErrUserNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	var adminName string
	privileges := false
	if admin != nil {
		adminName = admin.FullName().Name + " " + admin.FullName().Lastname
		if userID.Id == room.AdminID().Id {
			privileges = true
		}
	} else {
		adminName = "null"
	}

	var settingRoom *domain.SettingRoom
	var roomDetailedDTO *d.RoomDetailedDTO
	if room.IsFormal().IsFormal {
		settingRoom, err = h.GetSrByRoomIDUsecase.Execute(room.ID())

		if err != nil {
			if errors.Is(err, srerr.SettingRoomNotFoundError) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		roomDetailedDTO = &d.RoomDetailedDTO{
			ID:          room.ID().Id,
			RoomTitle:   room.Name().Name,
			AdminName:   adminName,
			Description: room.Description().Description,
			Code:        room.Code().Code,
			StartTime:   settingRoom.StartTime().DateTime,
			State:       room.State().CurrentState,
			Privileges:  privileges,
			Image:       room.Image().URL(),
		}
	} else {
		roomDetailedDTO = &d.RoomDetailedDTO{
			ID:          room.ID().Id,
			RoomTitle:   room.Name().Name,
			AdminName:   adminName,
			Description: room.Description().Description,
			Code:        room.Code().Code,
			State:       room.State().CurrentState,
			Privileges:  privileges,
			Image:       room.Image().URL(),
		}
	}

	response := map[string]interface{}{
		"success": "Ingreso a la sala exitoso",
		"room":    roomDetailedDTO,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *RoomEchoHandler) AddSingleUser(c echo.Context) error {
	var req d.AddSingleUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID, err := GetUserIDFromSession(c)

	if err != nil {
		return err
	}

	roomID, err := sv.NewID(req.RoomID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	err = h.AddSingleUserUsecase.Execute(req.UserData, *roomID, *userID)

	if err != nil {

		if errors.Is(err, rerr.ErrUserNotAdmin) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		} else if errors.Is(err, uerr.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		} else if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		} else if errors.Is(err, rerr.ErrAlreadyInWhitelist) {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

	}

	return c.JSON(http.StatusOK, map[string]string{"success": "usuario agregado a la sala exitosamente"})
}

func (h *RoomEchoHandler) Restore(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.RestoreUsecase.Execute(*id)
	if err != nil {
		if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "room not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "room restored succesfully"})

}

func (h *RoomEchoHandler) Update(c echo.Context) error {

	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room ID"})
	}

	var req d.RoomUpdate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	id, err := sv.NewID(roomID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	currentRoom, err := h.GetRoomByIDUsecase.Execute(*id)
	code := currentRoom.Code()

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid data"})
	}

	adminID, err := sv.NewID(currentRoom.AdminID().Id) // Usar el ID del admin actual o el que corresponda
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room AdminID"})
	}

	name, err := v.NewName(req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room name"})
	}

	description, err := v.NewDescription(req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room description"})
	}

	isFormal, err := v.NewIsFormal(req.IsFormal)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room isFormal"})
	}

	image, err := v.NewImage(req.Image)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid room Image"})
	}
	state, _ := v.NewState("created")

	room := d.NewRoom(
		id,
		*isFormal,
		&code,
		*name,
		adminID,
		*description,
		image,
		state,
	)

	userID, err := GetUserIDFromSession(c)

	if err != nil {
		return err
	}

	updatedRoom, err := h.UpdateRoomUsecase.Execute(room, *userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": err.Error()})
		}

		if errors.Is(rerr.ErrStateConstraint, err) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Devolver la respuesta con la sala actualizada
	roomDTO := d.RoomDTO{
		ID:          updatedRoom.ID().Id,
		IsFormal:    updatedRoom.IsFormal().IsFormal,
		Name:        updatedRoom.Name().Name,
		AdminID:     updatedRoom.AdminID().Id,
		Code:        updatedRoom.Code().Code,
		Description: updatedRoom.Description().Description,
		State:       room.State().CurrentState,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": "room updated successfully",
		"room":    roomDTO,
	})
}

func (r *RoomEchoHandler) RemoveFromWhitelistHandler(c echo.Context) error {

	var req d.RemoveFromWhitelistRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	adminId, _ := GetUserIDFromSession(c)

	roomId, err := sv.NewID(req.RoomId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	userId, err := sv.NewID(req.UserId)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = r.WhiteListRmUsecase.Execute(*roomId, *userId, *adminId)

	if err != nil {
		//varios tipos de errores
		if errors.Is(rerr.ErrRoomNotFound, err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		} else if errors.Is(uerr.ErrUserNotFound, err) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		} else if errors.Is(rerr.ErrUserNotAdmin, err) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"success": "room updated successfully"})
	return c.JSON(http.StatusOK, map[string]interface{}{"success": "user deleted from whitelist sucessfully"})
}

func (h *RoomEchoHandler) WsHandler(c echo.Context) error {

	id := c.Param("room_id")
	roomId, err := sv.NewID(id)
	if err != nil {
		return err
	}

	clientID, err := GetUserIDFromSession(c)
	if err != nil {
		return nil
	}

	//TODO: validar que sea el administrador

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	err = h.ManageWsUsecase.Execute(ws, *clientID, *roomId)
	if err != nil {
		ws.Close()
		log.Println(err.Error())
	}

	return nil
}

// En caso de devolver error lo hace en forma de response
func GetUserIDFromSession(c echo.Context) (*sv.ID, error) {
	// Obtener el user_id de la sesion
	userIDStr, ok := c.Get("user_id").(string)
	if !ok || userIDStr == "" {
		return nil, c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuario no autenticado"})
	}

	adminIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	adminID, err := sv.NewID(uint(adminIDUint))
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	return adminID, nil
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Permitir conexiones desde http://localhost:3000 (ajusta según tu frontend)
			origin := r.Header.Get("Origin")
			return origin == "http://localhost:4321"
		},
	}
)

func (h *RoomEchoHandler) HistoryRooms(c echo.Context) error {
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

	rooms, err := h.HistoryRoomsUsecase.Execute(*userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, rooms)
}
