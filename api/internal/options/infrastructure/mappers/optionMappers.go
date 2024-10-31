package mappers

import (
	"suffgo/internal/options/domain"
	v "suffgo/internal/options/domain/valueObjects"
	m "suffgo/internal/options/infrastructure/models"
	sv "suffgo/internal/shared/domain/valueObjects"
)

func DomainToModel(option *domain.Option) *m.Option {
	return &m.Option{
		ID:    option.ID().Id,
		Value: option.Value().Value,
	}
}

func ModelToDomain(optionModel *m.Option) (*domain.Option, error) {
	id, err := sv.NewID(optionModel.ID)
	if err != nil {
		return nil, err
	}

	value, err := v.NewValue(optionModel.Value)
	if err != nil {
		return nil, err
	}

	return domain.NewOption(id, *value), nil

}
