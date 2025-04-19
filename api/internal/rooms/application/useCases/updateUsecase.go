package usecases

import (
	"errors"
	"suffgo/internal/rooms/domain"
	e "suffgo/internal/rooms/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type UpdateRoomUsecase struct {
	repository domain.RoomRepository
}

func NewUpdateRoomUsecase(repository domain.RoomRepository) *UpdateRoomUsecase {
	return &UpdateRoomUsecase{
		repository: repository,
	}
}

func (u *UpdateRoomUsecase) Execute(room *domain.Room, userID sv.ID) (*domain.Room, error) {
	// Buscar la sala por ID

	existingRoom, err := u.repository.GetByID(room.ID())
	if err != nil {
		return nil, err
	}
	if existingRoom == nil {
		return nil, e.ErrRoomNotFound
	}

	//solo la puede editar si la sala no inicio
	if existingRoom.State().CurrentState !=  "created" {
		return nil, e.ErrStateConstraint
	} 

	if existingRoom.AdminID() != userID {
		return nil, errors.New("unauthorized")
	}

	// Guardar los cambios en el repositorio
	updatedRoom, err := u.repository.Update(room)
	if err != nil {
		return nil, err
	}

	return updatedRoom, nil
}
