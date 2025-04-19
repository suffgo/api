package usecases

import (
	"errors"
	d "suffgo/internal/proposals/domain"
	e "suffgo/internal/proposals/domain/errors"
	rd "suffgo/internal/rooms/domain"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type UpdateUsecase struct {
	repository      d.ProposalRepository
	roomRespository rd.RoomRepository
}

func NewUpdateProposalUsecase(repository d.ProposalRepository, roomRespository rd.RoomRepository) *UpdateUsecase {
	return &UpdateUsecase{
		repository:      repository,
		roomRespository: roomRespository,
	}
}

func (u *UpdateUsecase) Execute(proposal *d.Proposal, userID sv.ID) (*d.Proposal, error) {
	existingProposal, err := u.repository.GetById(proposal.ID())

	if err != nil {
		return nil, err
	}

	if existingProposal == nil {
		return nil, e.ErrPropNotFound
	}

	roomID, err := sv.NewID(proposal.RoomID().Id)
	if err != nil {
		return nil, err
	}

	rooom, err := u.roomRespository.GetByID(*roomID)
	if err != nil {
		return nil, err
	}

	if rooom.AdminID() != userID {
		return nil, errors.New("unauthorized")
	}

	updateProposal, err := u.repository.Update(proposal)

	if err != nil {
		return nil, err
	}

	return updateProposal, nil
}
