package socketStructs

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	userdom "suffgo/internal/users/domain"
)

type Client struct {
	conn   *websocket.Conn
	User   userdom.User
	lobby  *RoomLobby
	voted  bool
	egress chan Event //debido a que la conexion no soporta muchos mensajes al mismo tiempo, se utiliza este canal para que los mensajes lleguen uno a la vez
	done   chan struct{}
}

func NewClient(conn *websocket.Conn, user userdom.User) *Client {
	return &Client{
		conn:   conn,
		User:   user,
		voted:  false,
		egress: make(chan Event),
		done:   make(chan struct{}),
	}
}

func (c *Client) ReadMessages() {
	for {
		_, payload, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ws message: %v \n", err)
			}
			return
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling event : %v\n", err)
			break
		}

		if err := c.Lobby().routeEvent(request, c); err != nil {
			log.Println("error handling message: ", err)
			break
		}
	}

	c.lobby.removeClient(c)
}

func (c *Client) WriteMessages() {

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed", err)
				}
				break
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message: %v", err)
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
