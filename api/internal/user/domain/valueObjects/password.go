package valueobjects

type (
	Password struct {
		Password string
	}
)

func NewPassword(password string) (*Password, error) {
	return &Password{
		Password: password,
	}, nil
}