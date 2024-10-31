package valueobjects

import (
	"errors"
	"strings"
)

type (
	Dni struct {
		Dni string
	}
)

// falta devolver error opcionalmente en caso de error de validacion
func NewDni(dni string) (*Dni, error) {
	dni = strings.TrimSpace(dni)

	// validaciones re locas
	if dni == "" {
		return nil, errors.New("invalid dni")
	}

	if len(dni) < 8 {
		return nil, errors.New("invalid dni")
	}

	return &Dni{
		Dni: dni,
	}, nil
}
