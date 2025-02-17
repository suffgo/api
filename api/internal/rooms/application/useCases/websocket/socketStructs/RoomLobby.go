package socketStructs

import (
	"log"
	"suffgo/internal/rooms/domain"
	"sync"
)

type ClientList map[*Client]bool

type RoomLobby struct {
	sync.RWMutex
	clients ClientList //incluyendo al administrador
	admin   *Client
	room    *domain.Room

	usecases map[string]EventUsecase
}

func NewRoomLobby(admin *Client, room *domain.Room, roomRepo domain.RoomRepository) *RoomLobby {
	r := &RoomLobby{
		clients:  make(ClientList),
		admin:    admin,
		room:     room,
		usecases: make(map[string]EventUsecase),
	}

	r.initializeUsecases()

	return r
}

func (r *RoomLobby) initializeUsecases( /*repo domain.RoomRepository (por ahora no lo uso)*/ ) {
	r.usecases[EventSendMessage] = SendMessage
}

func SendMessage(event Event, c *Client) error {
	for client := range c.Lobby().Clients() {
		if client != c {
			client.egress <- event
		}
	}
	return nil
}

// me fijo que tipo de accion es y lo derivo al caso de uso correspondiente
func (r *RoomLobby) routeEvent(event Event, c *Client) error {
	if usecase, ok := r.usecases[event.Action]; ok {
		if err := usecase(event, c); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (r *RoomLobby) Admin() *Client {
	r.Lock()
	defer r.Unlock()

	return r.admin
}

func (r *RoomLobby) AddClient(client *Client) {
	r.Lock()
	defer r.Unlock()

	r.clients[client] = true //lo agrego a la lista de clientes conectados
	for user, conn := range r.clients {
		log.Printf("user %s; conn: %t", user.username, conn)
	}

}

func (r *RoomLobby) removeClient(client *Client) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.clients[client]; ok {
		client.conn.Close()
		delete(r.clients, client)
	}
}

func (r *RoomLobby) Clients() ClientList {
	return r.clients
}

func (r *RoomLobby) Room() *domain.Room {
	return r.room
}
