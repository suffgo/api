package valueobjects

type (
	UserDni struct {
		Dni string
	}
)

//falta devolver error opcionalmente en caso de error de validacion
func NewUserDni(dni string) *UserDni {
	
	// validaciones re locas

	return &UserDni{
		Dni: dni,
	}
}


