package usecases

import (
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
	rv "suffgo/internal/rooms/domain/valueObjects"
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


	code, err := s.roomGetByIDRepository.GetInviteCode(room.ID().Id)

	if err != nil {
		return nil, err
	}

	inviteCode, err  := rv.NewInviteCode(code)

	if err != nil {
		return nil, err
	}

	room.SetInviteCode(*inviteCode)

	return room, nil
}
