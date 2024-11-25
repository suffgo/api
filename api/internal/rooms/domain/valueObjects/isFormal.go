package valueobjects

type (
	IsFormal struct {
		IsFormal bool
	}
)

func NewIsFormal(isFormal bool) (*IsFormal, error) {
	return &IsFormal{
		IsFormal: isFormal,
	}, nil
}
