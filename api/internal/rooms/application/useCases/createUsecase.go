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

func (s *CreateUsecase) Execute(room domain.Room) error {

	err := s.repository.Save(room)
	if err != nil {
		return err
	}

	return nil

}
