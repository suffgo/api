package domain

import (
	v "suffgo/internal/settingsRoom/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"
	"time"
)

type (
	SettingRoom struct {
		id            *sv.ID
		privacy       v.Privacy
		proposalTimer v.ProposalTimer
		quorum        *v.Quorum
		startTime     *v.DateTime
		voterLimit    v.VoterLimit
		roomID        *sv.ID
	}

	SettingRoomDTO struct {
		ID            uint       `json:"id"`
		Privacy       bool       `json:"privacy"`
		ProposalTimer int        `json:"proposal_timer"`
		Quorum        *int       `json:"quorum"`
		StartTime     *time.Time `json:"start_time"`
		VoterLimit    int        `json:"voter_limit"`
		RoomID        uint       `json:"room_id"`
	}

	SettingRoomCreateRequest struct {
		Privacy       bool   `json:"privacy"`
		ProposalTimer int    `json:"proposal_timer"`
		Quorum        *int   `json:"quorum"`
		DateTime      *time.Time `json:"date_time"`
		VoterLimit    int    `json:"voter_limit"`
		RoomID        uint   `json:"room_id"`
	}
)

func NewSettingRoom(
	id *sv.ID,
	privacy v.Privacy,
	proposalTimer v.ProposalTimer,
	quorum v.Quorum,
	startTime v.DateTime,
	voterLimit v.VoterLimit,
	roomID *sv.ID,
) *SettingRoom {
	return &SettingRoom{
		id:            id,
		privacy:       privacy,
		proposalTimer: proposalTimer,
		quorum:        &quorum,
		startTime:     &startTime,
		voterLimit:    voterLimit,
		roomID:        roomID,
	}
}

func (s *SettingRoom) ID() sv.ID {
	return *s.id
}

func (s *SettingRoom) Privacy() v.Privacy {
	return s.privacy
}

func (s *SettingRoom) ProposalTimer() v.ProposalTimer {
	return s.proposalTimer
}

func (s *SettingRoom) Quorum() v.Quorum {
	return *s.quorum
}

func (s *SettingRoom) StartTime() v.DateTime {
	return *s.startTime
}

func (s *SettingRoom) VoterLimit() v.VoterLimit {
	return s.voterLimit
}

func (s *SettingRoom) RoomID() sv.ID {
	return *s.roomID
}
