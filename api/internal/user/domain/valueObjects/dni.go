package valueobjects

import "errors"

type (
	Dni struct {
		Dni string
	}
)

// falta devolver error opcionalmente en caso de error de validacion
func NewDni(dni string) (*Dni, error) {

	// validaciones re locas
	if dni == "" {
		return nil, errors.New("Invalid dni")
	}

	if len(dni) < 8 {
		return nil, errors.New("Invalid dni")
	}

	return &Dni{
		Dni: dni,
	}, nil
}
