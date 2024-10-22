package valueobjects

type (
	UserEmail struct {
		Email string
	}
)


func NewUserEmail(email string) *UserEmail {
	return &UserEmail{
		Email: email,
	}
}
