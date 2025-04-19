package usecases

import (
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type GetByIDUsecase struct {
	roomGetByIDRepository domain.RoomRepository
}

func NewGetByIDUsecase(repository domain.RoomRepository) *GetByIDUsecase {
	return &GetByIDUsecase{
		roomGetByIDRepository: repository,
	}
}

func (s *GetByIDUsecase) Execute(id sv.ID) (*domain.Room, error) {
	room, err := s.roomGetByIDRepository.GetByID(id)

	if err != nil {
		return nil, err
	}


	return room, nil
}
