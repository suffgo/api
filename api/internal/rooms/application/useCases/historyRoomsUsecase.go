package usecases

import (
	"suffgo/internal/rooms/domain"
	rv "suffgo/internal/rooms/domain/valueObjects"
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

	updatedRooms := make([]domain.Room, 0, len(rooms))

	for _, room := range rooms {
		code, err := s.historyRoomRepository.GetInviteCode(room.ID().Id)

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
