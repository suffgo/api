package infrastructure

import (
	"suffgo/cmd/database"
	"suffgo/internal/settingsRoom/domain"
	d "suffgo/internal/settingsRoom/domain"
	se "suffgo/internal/settingsRoom/domain/errors"
	"suffgo/internal/settingsRoom/infrastructure/mappers"
	m "suffgo/internal/settingsRoom/infrastructure/models"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type SettingRoomXormRepository struct {
	db database.Database
}

func NewSettingRoomXormRepository(db database.Database) *SettingRoomXormRepository {
	return &SettingRoomXormRepository{
		db: db,
	}
}

func (s *SettingRoomXormRepository) GetByID(id sv.ID) (*d.SettingRoom, error) {
	settingRoomModel := new(m.SettingsRoom)
	has, err := s.db.GetDb().ID(id.Id).Get(settingRoomModel)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, se.SettingRoomNotFoundError
	}
	userEnt, err := mappers.ModelToDomain(settingRoomModel)

	if err != nil {
		return nil, se.DataMappingError
	}
	return userEnt, nil
}

func (s *SettingRoomXormRepository) GetAll() ([]d.SettingRoom, error) {
	var settingsRoom []m.SettingsRoom
	err := s.db.GetDb().Find(&settingsRoom)
	if err != nil {
		return nil, err
	}

	var settingsRoomDomain []d.SettingRoom
	for _, settingRoom := range settingsRoom {
		settingRoomDomain, err := mappers.ModelToDomain(&settingRoom)

		if err != nil {
			return nil, err
		}

		settingsRoomDomain = append(settingsRoomDomain, *settingRoomDomain)
	}
	return settingsRoomDomain, nil
}

func (s *SettingRoomXormRepository) Delete(id sv.ID) error {

	affected, err := s.db.GetDb().ID(id.Id).Delete(&m.SettingsRoom{})
	if err != nil {
		return err
	}

	if affected == 0 {
		return se.SettingRoomNotFoundError
	}

	return nil
}

func (s *SettingRoomXormRepository) Save(settingRoom domain.SettingRoom) error {
	settingRoomModel := &m.SettingsRoom{
		Privacy:       settingRoom.Privacy().Privacy,
		ProposalTimer: settingRoom.ProposalTimer().ProposalTimer,
		Quorum:        settingRoom.Quorum().Quorum,
		Time:          settingRoom.TimeAndDate().Time,
		Date:          settingRoom.TimeAndDate().Date,
		VoterLimit:    settingRoom.VoterLimit().VoterLimit,
	}
	_, err := s.db.GetDb().Insert(settingRoomModel)
	if err != nil {
		return err
	}

	return nil
}
