package usecases

import (
	"suffgo/internal/rooms/domain"
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

	return createdRoom,nil

}
