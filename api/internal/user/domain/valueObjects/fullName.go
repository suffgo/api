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
		return nil, errors.New("Invalid name")
	}

	if lastname == "" {
		return nil, errors.New("Invalid lastname")
	}

	return &FullName{
		Name:     name,
		Lastname: lastname,
	}, nil
}
