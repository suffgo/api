package usecases

import (
	"suffgo/internal/rooms/domain"
	e "suffgo/internal/rooms/domain/errors"
)

type UpdateRoomUsecase struct {
	repository domain.RoomRepository
}

func NewUpdateRoomUsecase(repository domain.RoomRepository) *UpdateRoomUsecase {
	return &UpdateRoomUsecase{
		repository: repository,
	}
}

func (u *UpdateRoomUsecase) Execute(room *domain.Room) (*domain.Room, error) {
	// Buscar la sala por ID
	existingRoom, err := u.repository.GetByID(room.ID())
	if err != nil {
		return nil, err
	}
	if existingRoom == nil {
		return nil, e.ErrRoomNotFound
	}

	// Guardar los cambios en el repositorio
	updatedRoom, err := u.repository.Update(room)
	if err != nil {
		return nil, err
	}

	return updatedRoom, nil
}
