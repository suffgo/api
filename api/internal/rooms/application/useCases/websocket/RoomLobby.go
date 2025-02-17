package websocket

import (
	"sync"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type ClientList map[*Client]bool

type RoomLobby struct{
	sync.RWMutex
	clients ClientList //incluyendo al administrador
	admin *Client
	roomId sv.ID
}

func NewRoomLobby(admin *Client, roomId sv.ID) *RoomLobby{
	return &RoomLobby{
		clients: make(ClientList),
		admin: admin,
		roomId: roomId,
	}
}

func (r *RoomLobby) getAdmin() *Client{
	r.Lock()
	defer r.Unlock()

	return r.admin
}

func (r *RoomLobby) addClient(client *Client) {
	r.Lock()
	defer r.Unlock()

	r.clients[client] = true //lo agrego a la lista de clientes conectados
} 

func (r *RoomLobby) removeClient(client *Client) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.clients[client]; ok {
		client.conn.Close()
		delete(r.clients, client)
	}
}