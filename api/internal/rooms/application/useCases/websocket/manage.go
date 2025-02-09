package websocket

import (
	"errors"
	"fmt"
	"log"
	"suffgo/internal/rooms/domain"

	"github.com/gorilla/websocket"
)

type (
	ManageWsUsecase struct {
		repository domain.RoomRepository
	}
)

func NewManageWsUsecase(repository domain.RoomRepository) *ManageWsUsecase {
	return &ManageWsUsecase{
		repository: repository,
	}
}

func (s *ManageWsUsecase) Execute(ws *websocket.Conn, username string) error {

	_, msg, err := ws.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("Error inesperado de cierre de WebSocket: %v", err)
		}
		hub.Unregister <- ws
		return errors.New("Error websocket")
	}

	// Si no se pudo obtener el nombre, se verá como vacío
	response := fmt.Sprintf("%s: %s", username, msg)

	hub.Broadcast <- Message{Sender: ws, Data: []byte(response)}

	return nil
}
