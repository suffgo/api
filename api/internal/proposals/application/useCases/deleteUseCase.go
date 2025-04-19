package usecases

import (
	"errors"
	"suffgo/internal/proposals/domain"
	rd "suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	DeleteUseCase struct {
		Repository     domain.ProposalRepository
		RoomRepository rd.RoomRepository
	}
)

func NewDeleteUseCase(Proposalrepository domain.ProposalRepository, RoomRepository rd.RoomRepository) *DeleteUseCase {
	return &DeleteUseCase{
		Repository:     Proposalrepository,
		RoomRepository: RoomRepository,
	}
}

func (s *DeleteUseCase) Execute(id sv.ID, userID sv.ID) error {

	proposal, err := s.Repository.GetById(id)
	if err != nil {
		return err
	}

	roomID, err := sv.NewID(proposal.RoomID().Id)

	room, err := s.RoomRepository.GetByID(*roomID)
	if err != nil {
		return err
	}

	if room.AdminID() != userID {
		return errors.New("unauthorized")
	}

	err = s.Repository.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
