package valueobjects

type (
	Email struct {
		Email string
	}
)

func NewEmail(email string) (*Email, error) {
	return &Email{
		Email: email,
	}, nil
}
