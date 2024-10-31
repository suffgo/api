package infrastructure

import (
	"suffgo/cmd/database"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
	d "suffgo/internal/votes/domain"
	ve "suffgo/internal/votes/domain/errors"
	"suffgo/internal/votes/infrastructure/mappers"
	m "suffgo/internal/votes/infrastructure/models"
)

type VoteXormRepository struct {
	db database.Database
}

func NewVoteXormRepository(db database.Database) *VoteXormRepository {
	return &VoteXormRepository{
		db: db,
	}
}

func (s *VoteXormRepository) GetByID(id sv.ID) (*d.Vote, error) {
	voteModel := new(m.Vote)
	has, err := s.db.GetDb().ID(id.Id).Get(voteModel)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ve.VoteNotFoundError
	}

	voteEnt, err := mappers.ModelToDomain(voteModel)

	if err != nil {
		return nil, se.DataMappingError
	}

	return voteEnt, nil
}

func (s *VoteXormRepository) GetAll() ([]d.Vote, error) {
	var votes []m.Vote
	err := s.db.GetDb().Find(&votes)
	if err != nil {
		return nil, err
	}

	var votesDomain []d.Vote
	for _, vote := range votes {
		voteDomain, err := mappers.ModelToDomain(&vote)

		if err != nil {
			return nil, err
		}

		votesDomain = append(votesDomain, *voteDomain)
	}
	return votesDomain, nil

}

func (s *VoteXormRepository) Delete(id sv.ID) error {
	affected, err := s.db.GetDb().ID(id.Id).Delete(&m.Vote{})
	if err != nil {
		return err
	}

	if affected == 0 {
		return ve.VoteNotFoundError
	}
	return nil
}

func (s *VoteXormRepository) Save(vote d.Vote) error {
	voteModel := &m.Vote{
		UserID:   vote.UserID().Id,
		OptionID: vote.OptionID().Id,
	}

	_, err := s.db.GetDb().Insert(voteModel)
	if err != nil {
		return err
	}

	return nil
}
