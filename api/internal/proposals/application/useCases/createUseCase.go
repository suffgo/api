package usecases

import (
	"errors"
	"suffgo/internal/proposals/domain"
	roomDomain "suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type (
	CreateUsecase struct {
		proposalRepo domain.ProposalRepository
		roomRepo     roomDomain.RoomRepository
	}
)

func NewCreateUsecase(proposalRepo domain.ProposalRepository, roomRepo roomDomain.RoomRepository) *CreateUsecase {
	return &CreateUsecase{
		proposalRepo: proposalRepo,
		roomRepo:     roomRepo,
	}
}

func (s *CreateUsecase) Execute(proposal domain.Proposal, requesterUsr sv.ID) (*domain.Proposal, error) {

	//verificar existencia de sala
	room, err := s.roomRepo.GetByID(proposal.RoomID())
	if err != nil {
		return nil, err
	}

	if room == nil {
		return nil, errors.New("invalid room id")
	}

	//validar usuario creador de propuesta con admin de la sala
	if room.AdminID().Id != requesterUsr.Id {
		return nil, errors.New("operación no autorizada para este usuario")
	}

	createdProp, err := s.proposalRepo.Save(proposal)
	
	if err != nil {
		return nil, err
	}

	return createdProp, nil

}
