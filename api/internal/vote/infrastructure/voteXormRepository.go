package infrastructure

import (
	"suffgo/cmd/database"
	d "suffgo/internal/vote/domain"
	m "suffgo/internal/vote/infrastructure/models"
)

type VoteXormRepository struct {
	db database.Database
}

func NewVoteXormRepository(db database.Database) *VoteXormRepository {
	return &VoteXormRepository{
		db: db,
	}
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
