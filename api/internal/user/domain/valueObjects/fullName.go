package valueobjects

type (
	FullName struct {
		Name     string
		Lastname string
	}
)

func NewFullName(name, lastname string) (*FullName, error) {

	return &FullName{
		Name:     name,
		Lastname: lastname,
	}, nil
}
