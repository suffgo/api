package usecases

import (
	"suffgo/internal/settingsRoom/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type DeleteUsecase struct {
	settingRoomDeleteRepository domain.SettingRoomRepository
}

func NewDeleteUsecase(repository domain.SettingRoomRepository) *DeleteUsecase {
	return &DeleteUsecase{
		settingRoomDeleteRepository: repository,
	}
}

func (s *DeleteUsecase) Execute(id sv.ID) error {
	err := s.settingRoomDeleteRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
