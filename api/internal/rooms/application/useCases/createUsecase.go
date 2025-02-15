package usecases

import (
	"suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"

	"github.com/google/uuid"
)

type (
	CreateUsecase struct {
		roomRepository domain.RoomRepository
	}
)

func NewCreateUsecase(roomRepo domain.RoomRepository) *CreateUsecase {
	return &CreateUsecase{
		roomRepository: roomRepo,
	}
}

func (s *CreateUsecase) Execute(roomData domain.Room) (*domain.Room, error) {

	roomData.State().SetState("created") 

	createdRoom, err := s.roomRepository.Save(roomData)
	if err != nil {
		return nil, err
	}

	//Genero el codigo de invitacion

	inviteCode, err := v.NewInviteCode(uuid.New().String())

	if err != nil {
		return nil, err
	}

	createdRoom.SetInviteCode(*inviteCode)

	//guardo codigo

	err = s.roomRepository.SaveInviteCode(inviteCode.Code, createdRoom.ID().Id)

	//si es formal añado el admin a la whitelist

	if createdRoom.IsFormal().IsFormal {

		err = s.roomRepository.AddToWhitelist(createdRoom.ID(), createdRoom.AdminID())

		if err != nil {
			return nil, err
		}
	}

	return createdRoom, nil
}
