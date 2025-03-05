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
	clients   ClientList
	admin     *Client
	room      *domain.Room
	clientsmx sync.RWMutex
	proposals []propdom.Proposal
	propRepo  propdom.ProposalRepository
	roomRepo  domain.RoomRepository
	optRepo   optdom.OptionRepository
	voteRepo  votedom.VoteRepository
	usecases  map[string]EventUsecase
}

func NewRoomLobby(admin *Client, room *domain.Room, roomRepo domain.RoomRepository, propRepo propdom.ProposalRepository, optRepo optdom.OptionRepository, voteRepo votedom.VoteRepository) *RoomLobby {

	//error ya manejado anteriormente
	proposals, _ := propRepo.GetByRoom(room.ID())

	r := &RoomLobby{
		clients:   make(ClientList),
		admin:     admin,
		room:      room,
		usecases:  make(map[string]EventUsecase),
		proposals: proposals,
		roomRepo:  roomRepo,
		propRepo:  propRepo,
		optRepo:   optRepo,
		voteRepo:  voteRepo,
	}

	r.initializeUsecases()

	return r
}

func (r *RoomLobby) initializeUsecases() {
	r.usecases[EventSendMessage] = SendMessage
	r.usecases[EventStartVoting] = StartVoting
	r.usecases[EventVote] = ReceiveVote
}

func SendMessage(event Event, c *Client) error {
	for client := range c.Lobby().Clients() {
		if client != c {
			client.egress <- event
		}
	}
	return nil
}

func StartVoting(event Event, c *Client) error {
	log.Printf("room with id = %d has begun \n", c.Lobby().room.AdminID().Id)

	for _, prop := range c.Lobby().proposals {
		log.Println(prop)
	}

	if c.user.ID().Id != c.Lobby().Admin().user.ID().Id {

		errorEvent := Event{
			Action:  EventError,
			Payload: marshalOrPanic(ErrorEvent{Message: "You are not the admin"}),
		}

		c.egress <- errorEvent
		return nil
	}

	if len(c.Lobby().proposals) > 0 {
		proposal := c.Lobby().proposals[0]

		//obtengo opciones de la prop
		options, err := c.lobby.optRepo.GetByProposal(proposal.ID())
		if err != nil {
			errorEvent := Event{
				Action:  EventError,
				Payload: marshalOrPanic(ErrorEvent{Message: "error fetching options"}),
			}

			c.egress <- errorEvent

			return nil
		}

		var optionsValue []string
		for _, option := range options {
			optionsValue = append(optionsValue, option.Value().Value)
		}

		//todo lo necesario para poder votar
		proposalevt := ProposalEvent{
			ID:          proposal.ID().Id,
			Archive:     &proposal.Archive().Archive,
			Description: &proposal.Description().Description,
			Title:       proposal.Title().Title,
			Options:     optionsValue,
		}

		prop := Event{
			Action:  EventFirstProp,
			Payload: marshalOrPanic(proposalevt),
		}

		for client := range c.Lobby().Clients() {
			client.egress <- prop
		}
	}

	return nil
}

func ReceiveVote(event Event, c *Client) error {


	log.Println("holis")
	
	return nil
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
	var usernames []string
	for client := range r.clients {
		usernames = append(usernames, client.user.Username().Username)
	}

	// 2. Creamos el evento con la acción y el payload correspondiente.
	updateEventData := UpdateClientListEvent{
		Clients: usernames,
	}

	event := Event{
		Action:  EventUpdateClientList,
		Payload: marshalOrPanic(updateEventData), // ver la función marshalOrPanic abajo
	}

	for client := range r.clients {
		client.egress <- event
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
		log.Printf("user %s; conn: %t", user.user.Username().Username, conn)
	}

	r.broadcastClientList()

}

func (r *RoomLobby) removeClient(client *Client) {
	r.clientsmx.Lock()
	defer r.clientsmx.Unlock()

	if _, ok := r.clients[client]; ok {
		log.Printf("removing client %s", client.user.Username().Username)
		client.conn.Close()
		delete(r.clients, client)
	}

	r.broadcastClientList()
}

func (r *RoomLobby) Clients() ClientList {
	r.clientsmx.RLock()
	defer r.clientsmx.RUnlock()
	return r.clients
}

func (r *RoomLobby) Room() *domain.Room {
	return r.room
}
