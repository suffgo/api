package websocket

import (
	"log"
	"suffgo/internal/rooms/domain"

	"github.com/gorilla/websocket"
)

type StartWsUsecase struct {
	repository domain.RoomRepository
}

func NewStartWsUsecase(repository domain.RoomRepository) *StartWsUsecase {
	return &StartWsUsecase{
		repository: repository,
	}
}

func (s *StartWsUsecase) Execute(ws *websocket.Conn) error {
	// Registrar la conexión en el hub
	hub.Register <- ws

	// Enviar un mensaje de bienvenida
	if err := ws.WriteMessage(websocket.TextMessage, []byte("Bienvenido!!")); err != nil {
		log.Println("Error al enviar mensaje de bienvenida:", err)
		return err
	}

	return nil
}

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
