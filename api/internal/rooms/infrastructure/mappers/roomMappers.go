package mappers

import (
	"suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"
	m "suffgo/internal/rooms/infrastructure/models"
	sv "suffgo/internal/shared/domain/valueObjects"
)

func DomainToModel(room *domain.Room) *m.Room {
	return &m.Room{
		ID:         room.ID().Id,
		LinkInvite: ptr(room.LinkInvite().LinkInvite),
		IsFormal:   room.IsFormal().IsFormal,
		Name:       room.Name().Name,
		UserID:     room.AdminID().Id,
	}
}

func ModelToDomain(roomModel *m.Room) (*domain.Room, error) {
	id, err := sv.NewID(roomModel.ID)
	if err != nil {
		return nil, err
	}
	linkInvite, err := v.NewLinkInvite(*roomModel.LinkInvite)
	if err != nil {
		return nil, err
	}
	isFormal, err := v.NewIsFormal(roomModel.IsFormal)
	if err != nil {
		return nil, err
	}
	name, err := v.NewName(roomModel.Name)
	if err != nil {
		return nil, err
	}
	adminID, err := sv.NewID(roomModel.UserID)
	if err != nil {
		return nil, err
	}

	return domain.NewRoom(id, *linkInvite, *isFormal, *name, adminID), nil
}

func ptr(s string) *string {
	return &s
}
