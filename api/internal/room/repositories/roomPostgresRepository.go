package repositories

import (
	"suffgo/database"
	"suffgo/internal/room/entities"

	"github.com/labstack/gommon/log"
)

type roomPostgresRepository struct {
	db database.Database
}

func NewRoomPostgresRepository(db database.Database) RoomRepository {
	return &roomPostgresRepository{db: db}
}

func (r *roomPostgresRepository) InsertRoomData(in *entities.RoomDto) error {
	data := &entities.Room{
		LinkInvite: nil,
		IsFormal:   in.IsFormal,
		Name:       in.Name,
		AdminID:    in.AdminID,
	}

	result := r.db.GetDb().Create(data)

	if result.Error != nil {
		log.Errorf("InsertRoomData: %v", result.Error)
		return result.Error
	}

	log.Debugf("InsertRoomData: %v", result.RowsAffected)
	return nil
}

func (r *roomPostgresRepository) GetRoomByID(id int) (*entities.RoomDto, error) {
	var room entities.Room

	result := r.db.GetDb().Preload("AllowedUsers").First(&room, id)

	if result.Error != nil {
		return nil, result.Error
	}

	roomData := &entities.RoomDto{
		ID:           room.ID,
		IsFormal:     room.IsFormal,
		Name:         room.Name,
		AdminID:      room.AdminID,
	}

	return roomData, nil
}

func (r *roomPostgresRepository) DeleteRoom(id int) error {
	return nil
}

func (r *roomPostgresRepository) FetchAll() ([]entities.RoomDto, error) {
	return nil, nil
}
