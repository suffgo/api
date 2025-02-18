package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	"suffgo/internal/rooms/domain"
	roomerr "suffgo/internal/rooms/domain/errors"
	userdom "suffgo/internal/users/domain"

	"suffgo/internal/rooms/application/useCases/websocket/socketStructs"
	_ "suffgo/internal/rooms/application/useCases/websocket/socketStructs"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type ManageWsUsecase struct {
	rooms          map[sv.ID]*socketStructs.RoomLobby
	userRepository userdom.UserRepository
	roomRepository domain.RoomRepository
}

func NewManageWsUsecase(repo domain.RoomRepository, userRepo userdom.UserRepository) *ManageWsUsecase {

	return &ManageWsUsecase{
		roomRepository: repo,
		userRepository: userRepo,
		rooms:          make(map[sv.ID]*socketStructs.RoomLobby),
	}
}

func (s *ManageWsUsecase) Execute(ws *websocket.Conn, userId, roomId sv.ID) error {

	user, err := s.userRepository.GetByID(userId)
	if err != nil {
		return err
	}

	client := socketStructs.NewClient(ws, *user)
	//inicia la sala
	if s.rooms[roomId] == nil {
		room, err := s.roomRepository.GetByID(roomId)
		if err != nil {
			return err
		}
		// Verifica que room no sea nil
		if room == nil {
			return fmt.Errorf("room not found")
		}

		if user.ID().Id != room.AdminID().Id {
			return roomerr.ErrUserNotAdmin
		}

		// Usa roomId para almacenar la sala, de forma consistente.
		s.rooms[roomId] = socketStructs.NewRoomLobby(client, room, s.roomRepository)
		log.Printf("Sala iniciada con id = %d \n", room.ID().Id)
	}

	client.SetLobby(s.rooms[roomId])

	s.rooms[roomId].AddClient(client)

	go client.ReadMessages()
	go client.WriteMessages()

	return nil
}
