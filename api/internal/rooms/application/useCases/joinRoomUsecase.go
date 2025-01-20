package usecases

import (
	"errors"
	"fmt"
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type JoinRoomUsecase struct {
	joinRoomUsecaseRepository domain.RoomRepository
}

func NewJoinRoomUsecase(repository domain.RoomRepository) *JoinRoomUsecase {
	return &JoinRoomUsecase{
		joinRoomUsecaseRepository: repository,
	}
}

//Metodo importante, gestiona la union del usuario a la sala
func (s *JoinRoomUsecase) Execute(roomCode string) (*domain.Room, error) {
	//Obtener sala a traves de codigo
	roomID, err := s.joinRoomUsecaseRepository.GetRoomByCode(roomCode)
	
	fmt.Println("1..")
	//validar error, si es nulo, el codigo no es valido
	if err != nil {
		return nil, errors.New("codigo de sala invalido")
	}


	rID := sv.ID{Id: roomID}
	//obtener datos de sala
	room, err := s.joinRoomUsecaseRepository.GetByID(rID)
	fmt.Println("2..")
	if err != nil {
		return nil, errors.New("error al obtener la sala")
	}

	code := sv.Code{Code: roomCode}
	room.SetInviteCode(code)

	//Si la sala es formal verificar si el usuario puede unirse a la misma (tiene permiso, cantidad maxima, etc )

	fmt.Println("3..")
	return room, nil
}
