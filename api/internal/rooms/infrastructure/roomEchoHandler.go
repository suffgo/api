package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	r "suffgo/internal/rooms/application/useCases"
	addUsers "suffgo/internal/rooms/application/useCases/addUsers"

	d "suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"

	sv "suffgo/internal/shared/domain/valueObjects"

	se "suffgo/internal/shared/domain/errors"

	rerr "suffgo/internal/rooms/domain/errors"
	uerr "suffgo/internal/users/domain/errors"

	useruc "suffgo/internal/users/application/useCases"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)


type RoomEchoHandler struct {
	CreateRoomUsecase    *r.CreateUsecase
	DeleteRoomUsecase    *r.DeleteUsecase
	GetAllUsecase        *r.GetAllUsecase
	GetRoomByIDUsecase   *r.GetByIDUsecase
	GetByAdminUsecase    *r.GetByAdminUsecase
	RestoreUsecase       *r.RestoreUsecase
	JoinRoomUsecase      *r.JoinRoomUsecase
	AddSingleUSerUsecase *addUsers.AddSingleUserUsecase
	GetUserByIDUsecase   *useruc.GetByIDUsecase
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
) *RoomEchoHandler {
	return &RoomEchoHandler{
		CreateRoomUsecase:    creatUC,
		DeleteRoomUsecase:    deleteUC,
		GetAllUsecase:        getAllUC,
		GetRoomByIDUsecase:   getByIDUC,
		GetByAdminUsecase:    getByAdminUC,
		RestoreUsecase:       restoreUC,
		JoinRoomUsecase:      joinRoomUC,
		AddSingleUSerUsecase: addSingleUserUC,
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
	adminID, err := GetUserIDFromSession(c)

	if err != nil {
		return err
	}

	description, err := v.NewDescription(req.Description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	room := d.NewRoom(
		nil,
		*linkInvite,
		*isFormal,
		*name,
		adminID,
		*description,
	)

	createdRoom, err := h.CreateRoomUsecase.Execute(*room)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	roomDTO := &d.RoomDTO{
		ID:          createdRoom.ID().Id,
		LinkInvite:  createdRoom.LinkInvite().LinkInvite,
		IsFormal:    createdRoom.IsFormal().IsFormal,
		Name:        createdRoom.Name().Name,
		AdminID:     createdRoom.AdminID().Id,
		Description: createdRoom.Description().Description,
		RoomCode:    createdRoom.InviteCode().Code,
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
			ID:          room.ID().Id,
			LinkInvite:  room.LinkInvite().LinkInvite,
			IsFormal:    room.IsFormal().IsFormal,
			Name:        room.Name().Name,
			AdminID:     room.AdminID().Id,
			Description: room.Description().Description,
			RoomCode:    room.InviteCode().Code,
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
		ID:          room.ID().Id,
		LinkInvite:  room.LinkInvite().LinkInvite,
		IsFormal:    room.IsFormal().IsFormal,
		Name:        room.Name().Name,
		AdminID:     room.AdminID().Id,
		Description: room.Description().Description,
		RoomCode:    room.InviteCode().Code,
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
			RoomCode:   room.InviteCode().Code,
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
		if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	admin, err := userHandler.Execute(room.AdminID())

	if err != nil {
		if errors.Is(err, uerr.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	roomDetailedDTO := &d.RoomDetailedDTO{
		ID:          room.ID().Id,
		LinkInvite:  room.LinkInvite().LinkInvite,
		RoomTitle:   room.Name().Name,
		AdminName:   admin.FullName().Lastname + " " + admin.FullName().Name,
		Description: room.Description().Description,
		RoomCode:    room.InviteCode().Code,
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

	roomIDint, err := strconv.ParseUint(req.RoomID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	roomID, err := sv.NewID(uint(roomIDint))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	err = h.AddSingleUSerUsecase.Execute(req.UserData, *roomID, *userID)

	if err != nil {

		if errors.Is(err, rerr.ErrUserNotAdmin) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		if errors.Is(err, uerr.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}

		if errors.Is(err, rerr.ErrRoomNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}

		if errors.Is(err, rerr.ErrAlreadyInWhitelist) {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
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

	return c.JSON(http.StatusOK, map[string]string{"succes": "room restored succesfully"})

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

var hub = NewHub()

func init() {
	go hub.Run()
}

type Message struct {
	Sender *websocket.Conn
	Data   []byte
}

// Hub gestiona las conexiones y mensajes.
type Hub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan Message
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

// NewHub crea e inicializa un nuevo Hub.
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

// Run inicia el bucle principal del Hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("Cliente registrado. Total:", len(h.Clients))
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Close()
				log.Println("Cliente desregistrado. Total:", len(h.Clients))
			}
		case msg := <-h.Broadcast:
			// Enviar el mensaje a todos los clientes, excepto al remitente
			for client := range h.Clients {
				if client != msg.Sender {
					if err := client.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
						log.Println("Error al enviar mensaje a un cliente:", err)
						client.Close()
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}

func (h *RoomEchoHandler) WsHandler(c echo.Context) error {

	sess, err := session.Get("session", c)
	if err != nil {
		c.Logger().Error("Error al obtener la sesión:", err)
		return err
	}
	// Extraer el nombre desde la sesión
	username, _ := sess.Values["name"].(string)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("Error al actualizar a WebSocket:", err)
		return err
	}
	// Registrar la conexión en el hub
	hub.Register <- ws

	// Asegurarse de que la conexión se deba dar de baja cuando se cierre
	defer func() {
		hub.Unregister <- ws
		ws.Close()
	}()

	// Enviar un mensaje de bienvenida
	if err := ws.WriteMessage(websocket.TextMessage, []byte("Bienvenido!!")); err != nil {
		log.Println("Error al enviar mensaje de bienvenida:", err)
		return err
	}

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error inesperado de cierre de WebSocket: %v", err)
			}
			break
		}

		// Si no se pudo obtener el nombre, se verá como vacío
		response := fmt.Sprintf("%s: %s", username, msg)

		hub.Broadcast <- Message{Sender: ws, Data: []byte(response)}
	}

	return nil
}
