package valueobjects

import "errors"

type (
	FullName struct {
		Name     string
		Lastname string
	}
)

func NewFullName(name, lastname string) (*FullName, error) {
	if name == "" {
		return nil, errors.New("invalid name")
	}

	if lastname == "" {
		return nil, errors.New("invalid lastname")
	}

	return &FullName{
		Name:     name,
		Lastname: lastname,
	}, nil
}
