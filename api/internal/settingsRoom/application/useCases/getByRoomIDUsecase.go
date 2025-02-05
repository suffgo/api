package usecases

import (
	"suffgo/internal/settingsRoom/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type GetByRoomIDUsecase struct {
	GetByRoomIDUsecaseRepository domain.SettingRoomRepository
}

func NewGetByRoomID(repository domain.SettingRoomRepository) *GetByRoomIDUsecase {
	return &GetByRoomIDUsecase{
		GetByRoomIDUsecaseRepository: repository,
	}
}
func (s *GetByRoomIDUsecase) Execute(roomId sv.ID) (*domain.SettingRoom, error) {
	settingRoom, err := s.GetByRoomIDUsecaseRepository.GetByRoom(roomId)

	if err != nil {
		return nil, err
	}

	return settingRoom, nil
}
