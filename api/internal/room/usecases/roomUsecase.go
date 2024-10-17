package usecases

import (
	"suffgo/internal/room/entities"
	"suffgo/internal/room/models"
)

type RoomUsecase interface {
	RoomDataRegister(in *models.AddRoomData) error
	GetRoomByID(id string) (*entities.RoomDto, error)
	DeleteRoom(id string) error
	GetAll() ([]entities.RoomDto, error)
}
