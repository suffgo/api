package websocket

import (
	"suffgo/internal/rooms/domain"

	"github.com/gorilla/websocket"
)

type StartWsUsecase struct {
	repository domain.RoomRepository
}

func NewStartWsUsecase(repository domain.RoomRepository) *StartWsUsecase {
	return &StartWsUsecase{
		repository: repository,
	}
}

//Instancia la sala en RoomMap (solo el admin puede realizar esta accion)
func (s *StartWsUsecase) Execute(ws *websocket.Conn) error {

	return nil
}
