package usecases

import (
	sv "suffgo/internal/shared/domain/valueObjects"
	"suffgo/internal/rooms/domain"
)

type WhitelistRmUsecase struct {
	repo domain.RoomRepository
}

func NewWhitelistRmUsecase(repository domain.RoomRepository) *WhitelistRmUsecase {
	return &WhitelistRmUsecase{
		repo: repository,
	}
}

func (s *WhitelistRmUsecase) Execute(roomId, userId, adminId sv.ID) error {
	

	//validar sala

	//validar admin

	//validar usuario


	err := s.repo.RemoveFromWhitelist(roomId, userId)
	
	if err != nil {
		return err
	}
	
	return nil
}
