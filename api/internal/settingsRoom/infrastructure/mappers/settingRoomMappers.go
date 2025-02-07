package mappers

import (
	"suffgo/internal/settingsRoom/domain"
	v "suffgo/internal/settingsRoom/domain/valueObjects"
	m "suffgo/internal/settingsRoom/infrastructure/models"
	sv "suffgo/internal/shared/domain/valueObjects"
)

func DomainToModel(settingRoom *domain.SettingRoom) *m.SettingsRoom {
	return &m.SettingsRoom{
		ID:            settingRoom.ID().Id,
		Privacy:       settingRoom.Privacy().Privacy,
		ProposalTimer: settingRoom.ProposalTimer().ProposalTimer,
		Quorum:        settingRoom.Quorum().Quorum,
		StartTime:     settingRoom.StartTime().DateTime,
		VoterLimit:    settingRoom.VoterLimit().VoterLimit,
		RoomID:        settingRoom.RoomID().Id,
	}
}

func ModelToDomain(settingRoomModel *m.SettingsRoom) (*domain.SettingRoom, error) {
	id, err := sv.NewID(settingRoomModel.ID)
	if err != nil {
		return nil, err
	}

	privacy, err := v.NewPrivacy(settingRoomModel.Privacy)
	if err != nil {
		return nil, err
	}
	proposalTimer, err := v.NewProposalTimer(settingRoomModel.ProposalTimer)
	if err != nil {
		return nil, err
	}
	quorum, err := v.NewQuorum(settingRoomModel.Quorum)
	if err != nil {
		return nil, err
	}
	startTime, err := v.NewDateTime(settingRoomModel.StartTime)
	if err != nil {
		return nil, err
	}
	voterLimit, err := v.NewVoterLimit(settingRoomModel.VoterLimit)
	if err != nil {
		return nil, err
	}

	room, err := sv.NewID(settingRoomModel.RoomID)
	if err != nil {
		return nil, err
	}
	return domain.NewSettingRoom(id, *privacy, proposalTimer, *quorum, *startTime, voterLimit, room), nil
}
