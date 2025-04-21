package infrastructure

import (
	"errors"
	"log"
	"suffgo/cmd/database"
	"suffgo/internal/rooms/domain"
	d "suffgo/internal/rooms/domain"
	re "suffgo/internal/rooms/domain/errors"
	"suffgo/internal/rooms/infrastructure/mappers"
	m "suffgo/internal/rooms/infrastructure/models"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
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
		IsFormal:    room.IsFormal().IsFormal,
		Name:        room.Name().Name,
		AdminID:     room.AdminID().Id,
		Description: room.Description().Description,
		State:       room.State().CurrentState,
		Code:        room.Code().Code,
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

func (s *RoomXormRepository) GetRoomByCode(inviteCode string) (*d.Room, error) {
	//its only one room per code
	register := new(m.Room)
	has, err := s.db.GetDb().Where("code = ?", inviteCode).Get(register)

	log.Println("pichu")
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, re.ErrRoomNotFound
	}

	log.Println("pikachu")
	roomDom, err := mappers.ModelToDomain(register)
	if err != nil {
		return nil, se.ErrDataMap
	}

	log.Println("raichu")
	return roomDom, nil
}

// agrego un registro a user_room (para usuario registrado)
func (s *RoomXormRepository) AddToWhitelist(roomID sv.ID, userID sv.ID) error {

	reg := userRoomDom.UserRoom{
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
	var register []userRoomDom.UserRoom
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

func (s *RoomXormRepository) RestartRoom(roomId sv.ID) error {
	roomIDInt := int64(roomId.Value())

	_, err := s.db.GetDb().Exec(`
		DELETE FROM vote
		WHERE option_id IN (
			SELECT o.id
			FROM option o
			JOIN proposal p ON o.proposal_id = p.id
			WHERE p.room_id = ?
		);
	`, roomIDInt)

	if err != nil {
		return err
	}

	return nil
}

func (s *RoomXormRepository) HistoryRooms(userId sv.ID) ([]d.Room, error) {
	var roomModels []m.Room
	err := s.db.GetDb().SQL(`
        SELECT DISTINCT 
            r.id,
            r.is_formal,
            r.code,
            r.name,
            r.admin_id,
            r.description,
            r.state,
            COALESCE(r.image, '') AS image
        FROM vote v
        INNER JOIN option o ON v.option_id = o.id
        INNER JOIN proposal p ON o.proposal_id = p.id
        INNER JOIN room r ON p.room_id = r.id
        WHERE v.user_id = ?
        AND r.deleted_at IS NULL
        AND p.deleted_at IS NULL
    `, userId.Id).Find(&roomModels)

	if err != nil {
		return nil, err
	}

	var rooms []domain.Room
	for _, model := range roomModels {
		domainRoom, err := mappers.ModelToDomain(&model)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, *domainRoom)
	}

	return rooms, nil
}
