package usecases

import (
	"suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"
	domsettingroom "suffgo/internal/settingsRoom/domain"
	srv "suffgo/internal/settingsRoom/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"

	"github.com/google/uuid"
)

type (
	CreateUsecase struct {
		roomRepository  domain.RoomRepository
		settingRoomRepo domsettingroom.SettingRoomRepository
	}
)

func NewCreateUsecase(roomRepo domain.RoomRepository, srRepo domsettingroom.SettingRoomRepository) *CreateUsecase {
	return &CreateUsecase{
		roomRepository:  roomRepo,
		settingRoomRepo: srRepo,
	}
}

func (s *CreateUsecase) Execute(roomData domain.Room) (*domain.Room, error) {
	roomData.State().SetState("created")

	inviteCode, err := v.NewInviteCode(uuid.New().String()[:6])
	if err != nil {
		return nil, err
	}

	roomData.SetInviteCode(*inviteCode)

	createdRoom, err := s.roomRepository.Save(roomData)
	if err != nil {
		return nil, err
	}

	if createdRoom.IsFormal().IsFormal {
		err = s.roomRepository.AddToWhitelist(createdRoom.ID(), createdRoom.AdminID())
		if err != nil {
			return nil, err
		}
	}

	// Crear registro de settingRoom con valores por defecto
	settingRoom := generateDefaultRoomConfig(createdRoom.ID())
	err = s.settingRoomRepo.Save(settingRoom)
	if err != nil {
		return nil, err
	}

	return createdRoom, nil
}

func generateDefaultRoomConfig(roomId sv.ID) domsettingroom.SettingRoom {
	t := true
	zero := 0

	privacy, _ := srv.NewPrivacy(&t)
	proposalTimer, _ := srv.NewProposalTimer(30) //30 segundos por defecto
	quorum, _ := srv.NewQuorum(&zero)
	timeAndDate, _ := srv.NewDateTime(nil)
	voterLimit, _ := srv.NewVoterLimit(100)

	return *domsettingroom.NewSettingRoom(
		nil,
		*privacy,
		proposalTimer,
		*quorum,
		*timeAndDate,
		voterLimit,
		&roomId,
	)
}
