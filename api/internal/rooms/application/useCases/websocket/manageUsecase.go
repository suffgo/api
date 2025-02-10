package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"

	"suffgo/internal/rooms/domain"
)

// Message representa un mensaje que se difunde a los clientes conectados.
type Message struct {
	SenderID string
	Data   []byte
}

type Client struct {
    Username string
    Conn   *websocket.Conn
}

// Hub administra las conexiones y la difusión de mensajes en una sala.
type Hub struct {
    Clients    map[string]*Client // Clave: Username
    Broadcast  chan Message
    Register   chan *Client
    Unregister chan *Client
}

// NewHub crea e inicializa un nuevo Hub.
func NewHub() *Hub {
    return &Hub{
        Clients:    make(map[string]*Client),
        Broadcast:  make(chan Message),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),	
    }
}


// Run inicia el bucle principal del Hub, gestionando registros, desregistros y difusión.
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            // Si ya existe un cliente con ese UserID, cerramos la conexión anterior.
            if oldClient, exists := h.Clients[client.Username]; exists {
                log.Printf("El usuario %s ya estaba conectado, cerrando la conexión anterior.", client.Username)
                oldClient.Conn.Close()
            }
            h.Clients[client.Username] = client
            log.Printf("Cliente registrado: %s. Total clientes: %d", client.Username, len(h.Clients))

        case client := <-h.Unregister:
            if _, exists := h.Clients[client.Username]; exists {
                delete(h.Clients, client.Username)
                client.Conn.Close()
                log.Printf("Cliente desregistrado: %s. Total clientes: %d", client.Username, len(h.Clients))
            }

        case msg := <-h.Broadcast:
            // Se envía el mensaje a todos los clientes excepto al emisor.
            for uid, client := range h.Clients {
                if uid != msg.SenderID {
                    if err := client.Conn.WriteMessage(websocket.TextMessage, msg.Data); err != nil {
                        log.Printf("Error al enviar mensaje a %s: %v", uid, err)
                        client.Conn.Close()
                        delete(h.Clients, uid)
                    }
                }
            }
        }
    }
}



// RoomManager administra múltiples salas, cada una con su propio Hub.
type RoomManager struct {
	mu    sync.Mutex
	Rooms map[string]*Hub
}

// NewRoomManager crea e inicializa un RoomManager.
func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms: make(map[string]*Hub),
	}
}

// GetHub obtiene el Hub correspondiente a la sala dado su roomID. Si no existe, lo crea.
func (rm *RoomManager) GetHub(roomID string) *Hub {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	hub, exists := rm.Rooms[roomID]
	if !exists {
		hub = NewHub()
		rm.Rooms[roomID] = hub
		go hub.Run()
	}
	return hub
}

func (rm *RoomManager) InitializeHub(roomID string) *Hub{
	hub := NewHub()
	rm.Rooms[roomID] = hub
	go hub.Run()
	return hub
}

type ClientAction struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

type ManageWsUsecase struct {
	repository domain.RoomRepository
}

// NewManageWsUsecase crea una nueva instancia de ManageWsUsecase.
func NewManageWsUsecase(repo domain.RoomRepository) *ManageWsUsecase {
	return &ManageWsUsecase{
		repository: repo,
	}
}

var RoomMap = NewRoomManager()

func (s *ManageWsUsecase) Execute(ws *websocket.Conn, username, roomID string) error {
    // Obtener el hub correspondiente a la sala.
    hub := RoomMap.GetHub(roomID)

    // Leer el mensaje raw.
    _, rawMsg, err := ws.ReadMessage()
    if err != nil {
        if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
            log.Printf("Error inesperado al cerrar WebSocket: %v", err)
        }
        hub.Unregister <- &Client{Username: username, Conn: ws}
        return errors.New("error en la conexión websocket")
    }

    // Decodificar el JSON para identificar la acción.
    var clientAction ClientAction
    if err := json.Unmarshal(rawMsg, &clientAction); err != nil {
        log.Printf("Error al decodificar el mensaje JSON: %v", err)
        return err
    }

    switch clientAction.Action {
    case "send_message":
        var payload struct {
            Message string `json:"message"`
        }
        if err := json.Unmarshal(clientAction.Payload, &payload); err != nil {
            log.Printf("Error al decodificar el payload de send_message: %v", err)
            return err
        }
        response := fmt.Sprintf("%s: %s", username, payload.Message)
        log.Println(response)
        hub.Broadcast <- Message{SenderID: username, Data: []byte(response)}
    default:
        log.Printf("Acción desconocida: %s", clientAction.Action)
        errMsg := fmt.Sprintf("Acción %s no reconocida", clientAction.Action)
        if err := ws.WriteMessage(websocket.TextMessage, []byte(errMsg)); err != nil {
            log.Printf("Error al enviar mensaje de error: %v", err)
        }
    }
    return nil
}

