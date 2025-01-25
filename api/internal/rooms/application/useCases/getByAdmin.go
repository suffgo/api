package usecases

import (
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
	rv "suffgo/internal/rooms/domain/valueObjects"
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


	updatedRooms := make([]domain.Room, 0, len(rooms))

	for _, room := range rooms {
		code, err := s.roomGetByAdminRepository.GetInviteCode(room.ID().Id)

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
