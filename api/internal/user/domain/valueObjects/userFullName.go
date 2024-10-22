package valueobjects

type (
	UserFullName struct {
		Name     string
		Lastname string
	}
)

func NewUserFullName(name, lastname string) *UserFullName {

	return &UserFullName{
		Name:     name,
		Lastname: lastname,
	}
}