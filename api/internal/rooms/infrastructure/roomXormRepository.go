package infrastructure

import (
	"errors"
	"suffgo/cmd/database"
	d "suffgo/internal/rooms/domain"
	re "suffgo/internal/rooms/domain/errors"
	"suffgo/internal/rooms/infrastructure/mappers"
	m "suffgo/internal/rooms/infrastructure/models"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
	um "suffgo/internal/userRooms/infrastructure/models"
	userRoomDom "suffgo/internal/userRooms/infrastructure/models"
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
		LinkInvite:  ptr(room.LinkInvite().LinkInvite),
		IsFormal:    room.IsFormal().IsFormal,
		Name:        room.Name().Name,
		AdminID:     room.AdminID().Id,
		Description: room.Description().Description,
		State:       room.State().CurrentState,
		Image:       "",
	}

	if room.Image() != nil {
		roomModel.Image = room.Image().Image
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
	//its only one room per code
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

// agrego un registro a user_room (para usuario registrado)
func (s *RoomXormRepository) AddToWhitelist(roomID sv.ID, userID sv.ID) error {

	reg := um.UserRoom{
		UserID: userID.Id,
		RoomID: roomID.Id,
	}

	_, err := s.db.GetDb().Insert(&reg)

	if err != nil {
		return err
	}

	return nil
}

func (s *RoomXormRepository) UserInWhitelist(roomID sv.ID, userID sv.ID) (bool, error) {
	var register []um.UserRoom
	err := s.db.GetDb().Where("room_id = ? and user_id = ?", roomID.Id, userID.Id).Find(&register)

	if err != nil {
		return false, err
	}

	if register == nil {
		return false, nil
	}

	return true, nil
}

func (r *RoomXormRepository) Update(room *d.Room) (*d.Room, error) {
	roomID := room.ID().Id

	var existingRoom m.Room
	found, err := r.db.GetDb().ID(roomID).Get(&existingRoom)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("room not found")
	}

	updateRoom := mappers.DomainToModel(room)

	affected, err := r.db.GetDb().
		ID(roomID).
		Update(updateRoom)

	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, errors.New("no rows were updated")
	}

	updatedRoom, err := mappers.ModelToDomain(updateRoom)
	if err != nil {
		return nil, err
	}

	return updatedRoom, nil
}

func (s *RoomXormRepository) UpdateState(roomID sv.ID, state string) error {

	return nil
}

func (s *RoomXormRepository) RemoveFromWhitelist(roomId sv.ID, userId sv.ID) error {
	affected, err := s.db.GetDb().Where("room_id = ? AND user_id = ?", roomId.Id, userId.Id).Delete(&userRoomDom.UserRoom{})
	if err != nil {
		return err
	}

	if affected == 0 {
		return re.ErrRoomNotFound
	}

	return nil
}