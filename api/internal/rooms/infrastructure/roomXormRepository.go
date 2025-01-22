package infrastructure

import (
	"suffgo/cmd/database"
	d "suffgo/internal/rooms/domain"
	re "suffgo/internal/rooms/domain/errors"
	"suffgo/internal/rooms/infrastructure/mappers"
	m "suffgo/internal/rooms/infrastructure/models"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type RoomXormRepository struct {
	db database.Database
}

func NewRoomXormRepository(db database.Database) *RoomXormRepository {
	return &RoomXormRepository{
		db: db,
	}
}

func (s *RoomXormRepository) GetByID(id sv.ID) (*d.Room, error) {
	roomModel := new(m.Room)
	has, err := s.db.GetDb().ID(id.Id).Get(roomModel)

	if err != nil {
		return nil, err
	}
	if !has {
		return nil, re.ErrRoomNotFound
	}

	roomEnt, err := mappers.ModelToDomain(roomModel)

	if err != nil {
		return nil, se.ErrDataMap
	}

	return roomEnt, nil
}

func (s *RoomXormRepository) GetAll() ([]d.Room, error) {
	var rooms []m.Room
	err := s.db.GetDb().Where("deleted_at IS NULL").Find(&rooms)
	if err != nil {
		return nil, err
	}

	var roomsDomain []d.Room
	for _, rooms := range rooms {
		roomDomain, err := mappers.ModelToDomain(&rooms)

		if err != nil {
			return nil, err
		}

		roomsDomain = append(roomsDomain, *roomDomain)
	}
	return roomsDomain, nil
}

func (s *RoomXormRepository) Delete(id sv.ID) error {

	affected, err := s.db.GetDb().ID(id.Id).Delete(&m.Room{})
	if err != nil {
		return err
	}

	if affected == 0 {
		return re.ErrRoomNotFound
	}

	return nil
}

func (s *RoomXormRepository) Restore(roomID sv.ID) error {
	primitiveID := roomID.Value()

	user := &m.Room{DeletedAt: nil}

	affected, err := s.db.GetDb().Unscoped().ID(primitiveID).Cols("deleted_at").Update(user)
	if err != nil {
		return err
	}
	if affected == 0 {
		return re.ErrRoomNotFound
	}
	return err
}

func (r *RoomXormRepository) GetByAdminID(adminID sv.ID) ([]d.Room, error) {
	var rooms []m.Room
	err := r.db.GetDb().Where("admin_id = ?", adminID.Id).Find(&rooms)
	if err != nil {
		return nil, err
	}

	var roomsDomain []d.Room
	for _, room := range rooms {
		roomDomain, err := mappers.ModelToDomain(&room)

		if err != nil {
			return nil, se.ErrDataMap
		}

		roomsDomain = append(roomsDomain, *roomDomain)
	}
	return roomsDomain, nil
}

func (s *RoomXormRepository) Save(room d.Room) (*d.Room, error) {
	roomModel := &m.Room{
		LinkInvite: ptr(room.LinkInvite().LinkInvite),
		IsFormal:   room.IsFormal().IsFormal,
		Name:       room.Name().Name,
		AdminID:    room.AdminID().Id,
	}

	_, err := s.db.GetDb().Insert(roomModel)
	if err != nil {
		return nil, err
	}

	roomDom, err := mappers.ModelToDomain(roomModel)
	if err != nil {
		return nil, se.ErrDataMap
	}

	return roomDom, nil
}

func ptr(s string) *string {
	return &s
}

func (s *RoomXormRepository) SaveInviteCode(inviteCode string, roomID uint) error {
	inviteCodeModel := &m.InviteCode{
		RoomID: roomID,
		Code:   inviteCode,
	}

	_, err := s.db.GetDb().Insert(inviteCodeModel)
	if err != nil {
		return err
	}

	return nil
}

func (s *RoomXormRepository) GetInviteCode(roomID uint) (string, error) {
	var register []m.InviteCode
	err := s.db.GetDb().Where("room_id = ?", roomID).Find(&register)

	if err != nil {
		return "", err
	}

	return register[0].Code, nil
}

func (s *RoomXormRepository) GetRoomByCode(inviteCode string) (uint, error) {
	var register []m.InviteCode
	err := s.db.GetDb().Where("code = ?", inviteCode).Find(&register)

	if err != nil {
		return 0, err
	}

	if register == nil {
		return 0, re.ErrRoomNotFound
	}

	return register[0].RoomID, nil
}