package socketStructs

import (
	"encoding/json"
	"log"

	userdom "suffgo/internal/users/domain"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	User      userdom.User
	lobby     *RoomLobby
	voted     bool
	egress    chan Event
	done      chan struct{}
	errorSent chan struct{}
}

func NewClient(conn *websocket.Conn, user userdom.User) *Client {
	return &Client{
		conn:      conn,
		User:      user,
		voted:     false,
		egress:    make(chan Event),
		errorSent: make(chan struct{}),
		done:      make(chan struct{}),
	}
}

func (c *Client) ReadMessages() {
	for {
		_, payload, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ws message: %v \n", err)
			}
			c.lobby.removeClient(c)
			return
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling event : %v\n", err)
			c.lobby.removeClient(c)
			return
		}

		if err := c.Lobby().routeEvent(request, c); err != nil {
			log.Println("error handling message: ", err)
			c.lobby.removeClient(c)
			return
		}
	}
}

func (c *Client) WriteMessages() {

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message: %v", err)
			}

			var event Event
			if err := json.Unmarshal(data, &event); err != nil {
				log.Println("error unmarshalling message:", err)
			} else if event.Action == "kick_user" {
				c.errorSent <- struct{}{}
			}

		case <-c.done:
			close(c.egress)
			return
		}
	}
}

func (c *Client) Conn() *websocket.Conn {
	return c.conn
}

func (c *Client) SetConn(conn *websocket.Conn) {
	c.conn = conn
}

func (c *Client) Lobby() *RoomLobby {
	return c.lobby
}

func (c *Client) SetLobby(lobby *RoomLobby) {
	c.lobby = lobby
}
