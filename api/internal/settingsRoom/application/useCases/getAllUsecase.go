package usecases

import "suffgo/internal/settingsRoom/domain"

type (
	GetAllUsecase struct {
		getAllRepository domain.SettingRoomRepository
	}
)

func NewGetAllUsecase(repository domain.SettingRoomRepository) *GetAllUsecase {
	return &GetAllUsecase{
		getAllRepository: repository,
	}
}

func (s *GetAllUsecase) Execute() ([]domain.SettingRoom, error) {
	settingsRoom, err := s.getAllRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return settingsRoom, nil
}
