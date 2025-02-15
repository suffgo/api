package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
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

var RoomMap = NewRoomManager()

func (s *ManageWsUsecase) Execute(ws *websocket.Conn, username, roomID string, clientID sv.ID ) error {

    hub := RoomMap.GetHub(roomID)
    if hub == nil {
        ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4001, "Sala no existe"))
        return errors.New("no existe sala activa para la id dada.")
    }

    // Leer el mensaje raw.
    _, rawMsg, err := ws.ReadMessage()
    if err != nil {
        if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
            log.Printf("Error inesperado al cerrar WebSocket: %v", err)
        }
        hub.Unregister <- &Client{Username: username, Conn: ws}
        return errors.New("error en la conexión websocket")
    }

    // Decodificar el JSON para identificar la acción.
    var clientAction ClientAction
    if err := json.Unmarshal(rawMsg, &clientAction); err != nil {
        log.Printf("Error al decodificar el mensaje JSON: %v", err)
        return err
    }

	switch clientAction.Action {
	//en teoria esto no lo necesito, va a estar adentro del bucle de lectura
	case "send_message":
		var payload struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(clientAction.Payload, &payload); err != nil {
			log.Printf("Error al decodificar el payload de send_message: %v", err)
			return err
		}
		response := fmt.Sprintf("%s: %s", username, payload.Message)

		hub.Broadcast <- Message{SenderID: username, Data: []byte(response)}
	case "start_room": 
		//Obtengo la sala a partir del id

        roomIDobj, err := sv.NewID(roomID)
        if err != nil {
            return err
        }

		room, err := s.repository.GetByID(*roomIDobj)

        if err != nil {
            return err
        }

		if room.AdminID().Id != clientID.Id {
			return errors.ErrUnsupported
		}
	
		// Inicializar el hub de la sala
		hub := RoomMap.InitializeHub(room.ID().Id)
	
		// Registrar al administrador en el hub
		adminClient := &Client{Username: username, Conn: ws}
		hub.Register <- adminClient
	
		// Actualizar el estado de la sala a "online"
		s.repository.UpdateState(room.ID(), "online")
	
		// Notificar que el administrador se ha unido
		response := fmt.Sprintf("%s: se unió", username)
		hub.Broadcast <- Message{SenderID: username, Data: []byte(response)}


		//falta implementar go routine que contiene bucle de lectura 
		go hub.Run()
	case "join_room":
	default:
		log.Printf("Acción desconocida: %s", clientAction.Action)
		errMsg := fmt.Sprintf("Acción %s no reconocida", clientAction.Action)
		if err := ws.WriteMessage(websocket.TextMessage, []byte(errMsg)); err != nil {
			log.Printf("Error al enviar mensaje de error: %v", err)
		}

		
	}
	return nil
}
