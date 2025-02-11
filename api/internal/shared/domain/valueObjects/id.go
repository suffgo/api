package valueobjects

import (
	"errors"
	"strconv"
)

type (
	ID struct {
		Id uint
	}
)

func NewID(id interface{}) (*ID, error) {
	
	var uintID uint

	switch v := id.(type) {
	case uint:
		uintID = v
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, err
		}
		uintID = uint(parsed)
	case int:
		uintID = uint(v)
	default:
		return nil, errors.New("tipo no soportado: se espera uint o string")
	}

	return &ID{
		Id: uintID,
	}, nil

}

func (i *ID) Value() uint {
	return i.Id
}
