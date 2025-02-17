package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	username string
	lobby    *RoomLobby
	egress chan []byte //debido a que la conexion no soporta muchos mensajes al mismo tiempo, se utiliza este canal para que los mensajes lleguen uno a la vez
}

func NewClient(conn *websocket.Conn, username string) *Client {
	return &Client{
		conn:     conn,
		username: username,
		egress:   make(chan []byte),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.lobby.removeClient(c)
	}()

	for {
		msg, payload, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ws message: %v \n", err)
			}
			break
		}

		for wsclient := range c.lobby.clients {
			wsclient.egress <- payload	
		}

		log.Println(msg)
		log.Println(string(payload))
	}
}


func (c *Client) writeMessages() {
	defer func() {
		c.lobby.removeClient(c)
	}()
	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed", err)
				}
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("failed to send message")
			}
			log.Printf("%s says: %s \n", c.username, string(message))
		}
	}
}