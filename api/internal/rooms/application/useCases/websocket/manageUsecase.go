package websocket

import (
	"log"

	"github.com/gorilla/websocket"

	"suffgo/internal/rooms/domain"

	"suffgo/internal/rooms/application/useCases/websocket/socketStructs"
	_ "suffgo/internal/rooms/application/useCases/websocket/socketStructs"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type ManageWsUsecase struct {
	rooms          map[sv.ID]*socketStructs.RoomLobby
	roomRepository domain.RoomRepository
}

func NewManageWsUsecase(repo domain.RoomRepository) *ManageWsUsecase {

	return &ManageWsUsecase{
		roomRepository: repo,
		rooms:          make(map[sv.ID]*socketStructs.RoomLobby),
	}
}

func (s *ManageWsUsecase) Execute(ws *websocket.Conn, username string, roomId, clientID sv.ID) error {

	log.Printf("New connection %s \n", username)
	client := socketStructs.NewClient(ws, username)
	//inicia la sala
	if s.rooms[roomId] == nil {

		room, err := s.roomRepository.GetByID(roomId)

		if err != nil {
			return err
		}

		s.rooms[room.ID()] = socketStructs.NewRoomLobby(client, room, s.roomRepository)
		log.Printf("Sala iniciada con id = %d \n", room.ID().Id)
	}

	client.SetLobby(s.rooms[roomId])

	s.rooms[roomId].AddClient(client)

	go client.ReadMessages()
	go client.WriteMessages()

	return nil
}
