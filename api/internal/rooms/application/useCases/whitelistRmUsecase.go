package usecases

import (
	"suffgo/internal/rooms/domain"
	roomerr "suffgo/internal/rooms/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
	udom "suffgo/internal/users/domain"
	usererr "suffgo/internal/users/domain/errors"
)

type WhitelistRmUsecase struct {
	roomRep domain.RoomRepository
	userRep udom.UserRepository
}

func NewWhitelistRmUsecase(roomRepO domain.RoomRepository, userRepo udom.UserRepository) *WhitelistRmUsecase {
	return &WhitelistRmUsecase{
		roomRep: roomRepO,
		userRep: userRepo,
	}
}

func (s *WhitelistRmUsecase) Execute(roomId, userId, adminId sv.ID) error {

	//validar sala
	room, err := s.roomRep.GetByID(roomId)

	if err != nil {
		return roomerr.ErrRoomNotFound
	}

	//validar admin
	if room.AdminID().Id != adminId.Id {
		return roomerr.ErrUserNotAdmin
	}

	//validar usuario a ser borrado
	user, err := s.userRep.GetByID(userId)
	if err != nil {
		return err
	}

	if user == nil {
		return usererr.ErrUserNotFound
	}

	err = s.roomRep.RemoveFromWhitelist(roomId, userId)

	if err != nil {
		return err
	}

	return nil
}
