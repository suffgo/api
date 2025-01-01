package usecases

import (
	"suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"

	"github.com/google/uuid"
)

type (
	CreateUsecase struct {
		repository domain.RoomRepository
	}
)

func NewCreateUsecase(repository domain.RoomRepository) *CreateUsecase {
	return &CreateUsecase{
		repository: repository,
	}
}

func (s *CreateUsecase) Execute(roomData domain.Room) (*domain.Room, error) {

	createdRoom, err := s.repository.Save(roomData)
	if err != nil {
		return nil, err
	}

	//Genero y guardo el codigo de invitacion

	inviteCode, err := v.NewInviteCode(uuid.New().String())

	if err != nil {
		return nil, err
	}

	createdRoom.SetInviteCode(*inviteCode)

	

	return createdRoom,nil
}
