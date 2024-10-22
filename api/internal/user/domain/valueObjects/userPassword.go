package valueobjects

type (
	UserPassword struct {
		Password string
	}
)

func NewUserPassword(password string) *UserPassword {
	return &UserPassword{
		Password: password,
	}
}