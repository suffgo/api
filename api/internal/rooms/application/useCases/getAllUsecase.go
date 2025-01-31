package usecases

import (
	"suffgo/internal/rooms/domain"
	rv "suffgo/internal/rooms/domain/valueObjects"
)

type (
	GetAllUsecase struct {
		getAllRepository domain.RoomRepository
	}
)

func NewGetAllUsecase(repository domain.RoomRepository) *GetAllUsecase {
	return &GetAllUsecase{
		getAllRepository: repository,
	}
}

func (s *GetAllUsecase) Execute() ([]domain.Room, error) {

	rooms, err := s.getAllRepository.GetAll()

	if err != nil {
		return nil, err
	}

	updatedRooms := make([]domain.Room, 0, len(rooms))

	for _, room := range rooms {
		code, err := s.getAllRepository.GetInviteCode(room.ID().Id)

		if err != nil {
			return nil, err
		}

		inviteCode, err := rv.NewInviteCode(code)

		if err != nil {
			return nil, err
		}

		room.SetInviteCode(*inviteCode)

		updatedRooms = append(updatedRooms, room)
	}

	return updatedRooms, nil
}
