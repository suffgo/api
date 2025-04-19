package usecases

import (
	"suffgo/internal/proposals/domain"

	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	GetByRoomIDUsecase struct {
		getByRoomRepository domain.ProposalRepository
	}
)

func NewGetByRoomUsecase(repository domain.ProposalRepository) *GetByRoomIDUsecase {
	return &GetByRoomIDUsecase{
		getByRoomRepository: repository,
	}

}

func (s *GetByRoomIDUsecase) Execute(roomId *sv.ID) ([]domain.Proposal, error) {

	proposal, err := s.getByRoomRepository.GetByRoom(*roomId)

	if err != nil {
		return nil, err
	}

	return proposal, nil
}
