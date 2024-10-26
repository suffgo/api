package infrastructure

import (
	"suffgo/cmd/database"
	d "suffgo/internal/option/domain"
	oe "suffgo/internal/option/domain/errors"
	v "suffgo/internal/option/domain/valueObjects"
	"suffgo/internal/option/infrastructure/mappers"
	m "suffgo/internal/option/infrastructure/models"
	se "suffgo/internal/shared/domain/errors"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type OptionXormRepository struct {
	db database.Database
}

func NewOptionXormRepository(db database.Database) *OptionXormRepository {
	return &OptionXormRepository{
		db: db,
	}
}

func (s *OptionXormRepository) GetByID(id sv.ID) (*d.Option, error) {
	optionModel := new(m.Option)
	has, err := s.db.GetDb().ID(id.Id).Get(optionModel)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, oe.OptionNotFoundError
	}

	userEnt, err := mappers.ModelToDomain(optionModel)

	if err != nil {
		return nil, se.DataMappingError
	}

	return userEnt, nil
}

func (s *OptionXormRepository) GetByValue(value v.Value) (*d.Option, error) {
	optionModel := new(m.Option)
	has, err := s.db.GetDb().Where("value = ?", value.Value).Get(optionModel)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	optionEnt, err := mappers.ModelToDomain(optionModel)

	if err != nil {
		return nil, se.DataMappingError
	}

	return optionEnt, nil
}

func (s *OptionXormRepository) GetAll() ([]d.Option, error) {
	var options []m.Option
	err := s.db.GetDb().Find(&options)
	if err != nil {
		return nil, err
	}

	var optionsDomain []d.Option
	for _, option := range options {
		optionDomain, err := mappers.ModelToDomain(&option)

		if err != nil {
			return nil, err
		}
		optionsDomain = append(optionsDomain, *optionDomain)
	}
	return optionsDomain, nil
}

func (s *OptionXormRepository) Delete(id sv.ID) error {

	affected, err := s.db.GetDb().ID(id.Id).Delete(&m.Option{})
	if err != nil {
		return err
	}

	if affected == 0 {
		return oe.OptionNotFoundError
	}

	return nil
}

func (s *OptionXormRepository) Save(option d.Option) error {
	optionModel := &m.Option{
		Value: option.Value().Value,
	}

	_, err := s.db.GetDb().Insert(optionModel)
	if err != nil {
		return err
	}

	return nil
}
