package socketStructs

import (
	"encoding/json"
	"log"

	optdom "suffgo/internal/options/domain"
	propdom "suffgo/internal/proposals/domain"
	"suffgo/internal/rooms/domain"

	votedom "suffgo/internal/votes/domain"
	"sync"
)

type ClientList map[*Client]bool

type RoomLobby struct {
	sync.RWMutex
	clientsmx      sync.RWMutex
	votesProcesing chan struct{}
	Empty          chan struct{}

	clients      ClientList
	admin        *Client
	room         *domain.Room
	proposals    []propdom.Proposal
	propRepo     propdom.ProposalRepository
	roomRepo     domain.RoomRepository
	optRepo      optdom.OptionRepository
	voteRepo     votedom.VoteRepository
	usecases     map[string]EventUsecase
	results      map[*Client]votedom.Vote
	nextProposal int
}

func NewRoomLobby(admin *Client, room *domain.Room, roomRepo domain.RoomRepository, propRepo propdom.ProposalRepository, optRepo optdom.OptionRepository, voteRepo votedom.VoteRepository) *RoomLobby {

	//error ya manejado anteriormente
	proposals, _ := propRepo.GetByRoom(room.ID())

	r := &RoomLobby{
		clients:        make(ClientList),
		admin:          admin,
		room:           room,
		usecases:       make(map[string]EventUsecase),
		proposals:      proposals,
		roomRepo:       roomRepo,
		propRepo:       propRepo,
		optRepo:        optRepo,
		voteRepo:       voteRepo,
		results:        make(map[*Client]votedom.Vote),
		votesProcesing: make(chan struct{}, 1),
		nextProposal:   0,
		Empty:          make(chan struct{}, 1),
	}

	r.initializeUsecases()
	r.votesProcesing <- struct{}{}

	return r
}

func (r *RoomLobby) initializeUsecases() {
	r.usecases[EventSendMessage] = SendMessage
	r.usecases[EventStartVoting] = StartVoting
	r.usecases[EventVote] = ReceiveVote
	r.usecases[EventResults] = SendResults
	r.usecases[EventNextProp] = NextProposal
	r.usecases[EventKickUser] = KickUser
}

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

func (r *RoomLobby) broadcastClientList() {
	// 1. Recorremos los clientes activos para obtener sus nombres (o información requerida).
	var clients []ClientData
	for client := range r.clients {
		clientData := ClientData{
			ID:       client.User.ID().Id,
			Name:     client.User.FullName().Name,
			Lastname: client.User.FullName().Lastname,
			Username: client.User.Username().Username,
			Email:    client.User.Email().Email,
			Voted:    client.voted,
			Image:    client.User.Image().URL(),
		}

		clients = append(clients, clientData)
	}

	// 2. Creamos el evento con la acción y el payload correspondiente.
	updateEventData := UpdateClientListEvent{
		Clients: clients,
	}

	event := Event{
		Action:  EventUpdateClientList,
		Payload: marshalOrPanic(updateEventData),
	}

	for client := range r.clients {
		if r.clients[client] {
			client.egress <- event
		}

	}
}

func marshalOrPanic(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Panicln("error marshalling:", err)
	}
	return data
}

func (r *RoomLobby) AddClient(client *Client) {
	r.clientsmx.Lock()
	defer r.clientsmx.Unlock()

	r.clients[client] = true //lo agrego a la lista de clientes conectados
	for user, conn := range r.clients {
		log.Printf("user %s; conn: %t", user.User.Username().Username, conn)
	}

	r.broadcastClientList()

}

func (r *RoomLobby) removeClient(client *Client) {
	r.clientsmx.Lock()
	if _, ok := r.clients[client]; ok {
		client.conn.Close()
		r.clients[client] = false //estado desconectado
		delete(r.clients, client)
		close(client.done)
		//delete(r.clients, client)
	}
	r.clientsmx.Unlock()

	r.broadcastClientList()

	state := r.room.State().CurrentState
	if len(r.clients) == 0 && (state == "in progress" || state == "online") {
		r.room.State().SetState("created")
		_, err := r.roomRepo.Update(r.room)
		if err != nil {
			return
		}
		r.Empty <- struct{}{}

	}
}

func (r *RoomLobby) Clients() ClientList {
	r.clientsmx.RLock()
	defer r.clientsmx.RUnlock()
	return r.clients
}

func (r *RoomLobby) Room() *domain.Room {
	return r.room
}
