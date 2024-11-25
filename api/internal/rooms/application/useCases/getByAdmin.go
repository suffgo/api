package usecases

import (
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type GetByAdminUsecase struct {
	roomGetByAdminRepository domain.RoomRepository
}

func NewGetByAdminUsecase(repository domain.RoomRepository) *GetByAdminUsecase {
	return &GetByAdminUsecase{
		roomGetByAdminRepository: repository,
	}
}

func (s *GetByAdminUsecase) Execute(adminID sv.ID) ([]domain.Room, error) {
	rooms, err := s.roomGetByAdminRepository.GetByAdminID(adminID)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
