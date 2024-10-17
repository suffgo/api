package repositories

import "suffgo/internal/room/entities"

type RoomRepository interface {
	InsertRoomData(in *entities.RoomDto) error
	GetRoomByID(id int) (*entities.RoomDto, error)
	DeleteRoom(id int) error
	FetchAll() ([]entities.RoomDto, error)
}
