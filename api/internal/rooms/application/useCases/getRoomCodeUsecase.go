package usecases

import (
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type GetRoomCodeUsecase struct {
	roomCodeRepository domain.RoomRepository
}

func NewGetRoomCodeUsecase(repository domain.RoomRepository) *GetRoomCodeUsecase {
	return &GetRoomCodeUsecase{
		roomCodeRepository: repository,
	}
}

func (s *GetRoomCodeUsecase) Execute(id sv.ID) error {
	return nil
}