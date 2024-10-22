package valueobjects

type (
	UserName struct {
		Username string
	}
)

func NewUserName(username string) (*UserName, error) {
	return &UserName{
		Username: username,
	}, nil
}
