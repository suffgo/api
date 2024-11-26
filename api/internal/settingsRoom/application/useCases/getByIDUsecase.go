package usecases

import (
	"suffgo/internal/settingsRoom/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type GetByIDUsecase struct {
	settingRoomGetByIDRepository domain.SettingRoomRepository
}

func NewGetByIDUsecase(repository domain.SettingRoomRepository) *GetByIDUsecase {
	return &GetByIDUsecase{
		settingRoomGetByIDRepository: repository,
	}
}
func (s *GetByIDUsecase) Execute(id sv.ID) (*domain.SettingRoom, error) {
	settingRoom, err := s.settingRoomGetByIDRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return settingRoom, nil
}
