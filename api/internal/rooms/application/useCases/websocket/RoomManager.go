package websocket

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientAction struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

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

func (rm *RoomManager) GetHub(roomID string) *Hub {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	hub, exists := rm.Rooms[roomID]
	if !exists {
		return nil
	}
	return hub
}

// Solo lo debe poder usar el administrador
func (rm *RoomManager) InitializeHub(roomID uint) *Hub {
	roomIDStr := strconv.FormatUint(uint64(roomID), 10)
	hub := NewHub()
	rm.Rooms[roomIDStr] = hub
	return hub
}


// Message representa un mensaje que se difunde a los clientes conectados.
type Message struct {
	SenderID string
	Data     []byte
}

type Client struct {
	Username string
	Conn     *websocket.Conn
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