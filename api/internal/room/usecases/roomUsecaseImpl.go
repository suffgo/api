package usecases

import (
	"suffgo/internal/room/entities"
	"suffgo/internal/room/models"
	"suffgo/internal/room/repositories"
)

type roomUsecaseImpl struct {
	roomRepository repositories.RoomRepository
}

func NewRoomUsecaseImpl(roomRepository repositories.RoomRepository) RoomUsecase {
	return &roomUsecaseImpl{
		roomRepository: roomRepository,
	}
}

func (r *roomUsecaseImpl) RoomDataRegister(in *models.AddRoomData) error {
	return nil
}

func (r *roomUsecaseImpl) GetRoomByID(id string) (*entities.RoomDto, error) {
	return nil, nil
}

func (r *roomUsecaseImpl) DeleteRoom(id string) error {
	return nil
}

func (r *roomUsecaseImpl) GetAll() ([]entities.RoomDto, error) {
	return nil, nil
}
