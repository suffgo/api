package websocket

import (
	"log"

	"github.com/gorilla/websocket"

	"suffgo/internal/rooms/domain"

	sv "suffgo/internal/shared/domain/valueObjects"
)

type ManageWsUsecase struct {
	repository domain.RoomRepository
}

// NewManageWsUsecase crea una nueva instancia de ManageWsUsecase.
func NewManageWsUsecase(repo domain.RoomRepository) *ManageWsUsecase {
	return &ManageWsUsecase{
		repository: repo,
	}
}

//clave = id y valor = sala
var rooms map[sv.ID]*RoomLobby

//client = user
func (s *ManageWsUsecase) Execute(ws *websocket.Conn, username string, primitiveRoomId, clientID sv.ID) error {
	
	if rooms == nil {
		rooms = make(map[sv.ID]*RoomLobby)
	}

	log.Printf("New connection %s \n", username)

	client := NewClient(ws, username)
	//inicia la sala
	if rooms[primitiveRoomId] == nil {
		rooms[primitiveRoomId] = NewRoomLobby(client, primitiveRoomId)
	}
	
	//asigno el lobby al cliente
	client.lobby = rooms[primitiveRoomId]
	log.Printf("Sala iniciada con id = %d \n", primitiveRoomId.Id)
	
	//agrego al admin
	rooms[primitiveRoomId].addClient(client)

	go client.readMessages()
	go client.writeMessages()

	return nil
}