package valueobjects

type (
	UserUserName struct {
		Username string
	}
)

func NewUserUserName(username string) *UserUserName {
	return &UserUserName{
		Username: username,
	}
}