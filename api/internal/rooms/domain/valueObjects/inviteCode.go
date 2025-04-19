package valueobjects


type (
	InviteCode struct {
		Code string
	}
)

func NewInviteCode(code string) (*InviteCode, error) {
	return &InviteCode{
		Code: code,
	}, nil
}
