package valueobjects

type (
	Dni struct {
		Dni string
	}
)

// falta devolver error opcionalmente en caso de error de validacion
func NewDni(dni string) (*Dni, error) {

	// validaciones re locas

	return &Dni{
		Dni: dni,
	}, nil
}
