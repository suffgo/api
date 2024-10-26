package domain

import (
	v "suffgo/internal/option/domain/valueObjects"
	sv "suffgo/internal/shared/domain/valueObjects"
)

type OptionRepository interface {
	GetByID(id sv.ID) (*Option, error)
	GetAll() ([]Option, error)
	GetByValue(value v.Value) (*Option, error)
	Delete(id sv.ID) error
	Save(option Option) error
}
