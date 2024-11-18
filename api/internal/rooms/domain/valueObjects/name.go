package valueobjects

import "errors"

type (
	Name struct {
		Name string
	}
)

func NewName(name string) (*Name, error) {
	if name == "" {
		return nil, errors.New("invalid name")
	}

	return &Name{
		Name: name,
	}, nil
}
