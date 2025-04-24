package websocket

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"

	optdom "suffgo/internal/options/domain"
	propdom "suffgo/internal/proposals/domain"
	"suffgo/internal/rooms/domain"
	roomerr "suffgo/internal/rooms/domain/errors"
	userdom "suffgo/internal/users/domain"
	votedom "suffgo/internal/votes/domain"

	"suffgo/internal/rooms/application/useCases/websocket/socketStructs"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type ManageWsUsecase struct {
	rooms        map[sv.ID]*socketStructs.RoomLobby
	userRepo     userdom.UserRepository
	roomRepo     domain.RoomRepository
	proposalRepo propdom.ProposalRepository
	optionsRepo  optdom.OptionRepository
	voteRepo     votedom.VoteRepository
}

func NewManageWsUsecase(
	repo domain.RoomRepository,
	userRepo userdom.UserRepository,
	proposalRepo propdom.ProposalRepository,
	optionsRepo optdom.OptionRepository,
	votesRepo votedom.VoteRepository,
) *ManageWsUsecase {

	return &ManageWsUsecase{
		roomRepo:     repo,
		userRepo:     userRepo,
		proposalRepo: proposalRepo,
		optionsRepo:  optionsRepo,
		voteRepo:     votesRepo,
		rooms:        make(map[sv.ID]*socketStructs.RoomLobby),
	}
}

func (s *ManageWsUsecase) Execute(ws *websocket.Conn, userId, roomId sv.ID) error {

	user, err := s.userRepo.GetByID(userId)
	if err != nil {
		return err
	}
	var client *socketStructs.Client
	reconnect := false
	if s.rooms[roomId] != nil {
		for cKey := range s.rooms[roomId].Clients() {
			if cKey.User.ID().Id == user.ID().Id {
				// Ya está conectado => rechazamos la nueva conexión
				ws.WriteControl(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(4002, "Ya estas conectado a la sala"),
					time.Now().Add(time.Second),
				)
				return nil
			}
		}
	}

	if !reconnect {
		client = socketStructs.NewClient(ws, *user)
	}

	if s.rooms[roomId] == nil {
		room, err := s.roomRepo.GetByID(roomId)
		if err != nil {
			return err
		}
		if room == nil {
			return fmt.Errorf("room not found")
		}
		
		if room.State().CurrentState == "finished" {
			return nil
		}
		
		if user.ID().Id != room.AdminID().Id {
			return roomerr.ErrUserNotAdmin
		}

		room.State().SetState("online")
		updatedroom, err := s.roomRepo.Update(room)

		if err != nil {
			return err
		}

		s.rooms[roomId] = socketStructs.NewRoomLobby(
			client,
			updatedroom,
			s.roomRepo,
			s.proposalRepo,
			s.optionsRepo,
			s.voteRepo,
		)

		go s.OnEmpty(s.rooms[roomId])
	}

	client.SetLobby(s.rooms[roomId])

	go client.ReadMessages()
	go client.WriteMessages()

	s.rooms[roomId].AddClient(client)

	return nil
}

func (s *ManageWsUsecase) OnEmpty(room *socketStructs.RoomLobby) {
	<-room.Empty
	delete(s.rooms, room.Room().ID())
	log.Printf("Room instance cleared id = %d \n", room.Room().ID().Id)
}
