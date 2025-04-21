package usecases

import (
	"suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type HistoryRooms struct {
	historyRoomRepository domain.RoomRepository
}

func NewHistoryRoomsUsecase(repository domain.RoomRepository) *HistoryRooms {
	return &HistoryRooms{
		historyRoomRepository: repository,
	}
}

func (s *HistoryRooms) Execute(userID sv.ID) ([]domain.Room, error) {
	rooms, err := s.historyRoomRepository.HistoryRooms(userID)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
