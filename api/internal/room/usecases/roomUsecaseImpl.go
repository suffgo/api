package usecases

import (
	"strconv"
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
	insertRoomData := &entities.RoomDto{
		LinkInvite: in.LinkInvite,
		IsFormal:   in.IsFormal,
		Name:       in.Name,
		AdminID:    in.AdminID,
	}

	if err := r.roomRepository.InsertRoomData(insertRoomData); err != nil {
		return err
	}

	return nil
}

func (r *roomUsecaseImpl) GetRoomByID(id string) (*entities.RoomDto, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	user, err := r.roomRepository.GetRoomByID(idInt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *roomUsecaseImpl) DeleteRoom(id string) error {
	return nil
}

func (r *roomUsecaseImpl) GetAll() ([]entities.RoomDto, error) {
	return nil, nil
}
