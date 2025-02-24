package mappers

import (
	"suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"
	m "suffgo/internal/rooms/infrastructure/models"
	sv "suffgo/internal/shared/domain/valueObjects"
)

func DomainToModel(room *domain.Room) *m.Room {
	image := ""
	if room.Image() != nil {
		image = room.Image().Path()
	}
	return &m.Room{
		ID:          room.ID().Id,
		LinkInvite:  ptr(room.LinkInvite().LinkInvite),
		IsFormal:    room.IsFormal().IsFormal,
		Name:        room.Name().Name,
		Description: room.Description().Description,
		AdminID:     room.AdminID().Id,
		State:       room.State().CurrentState,
		Image:       image,
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
	adminID, err := sv.NewID(roomModel.AdminID)
	if err != nil {
		return nil, err
	}

	description, err := v.NewDescription(roomModel.Description)
	if err != nil {
		return nil, err
	}

	var image *v.Image
	if roomModel.Image != "" {
		image, err = v.NewImage(roomModel.Image)
		if err != nil {
			return nil, err
		}
	}

	return domain.NewRoom(id, *linkInvite, *isFormal, *name, adminID, *description, image), nil
}

func ptr(s string) *string {
	return &s
}
