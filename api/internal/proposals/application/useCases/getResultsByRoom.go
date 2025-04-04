package usecases

import (
	"suffgo/internal/proposals/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	GetResultsByRoomUsecase struct {
		getResultsByRoomRepository domain.ProposalRepository
	}
)

func NewGetResultsByRoomUsecase(repository domain.ProposalRepository) *GetResultsByRoomUsecase {
	return &GetResultsByRoomUsecase{
		getResultsByRoomRepository: repository,
	}
}

func (s *GetResultsByRoomUsecase) Execute(roomId *sv.ID) ([]domain.ProposalResults, error) {
	proposal, err := s.getResultsByRoomRepository.GetResultsByRoom(*roomId)

	if err != nil {
		return nil, err
	}

	return proposal, nil
}
