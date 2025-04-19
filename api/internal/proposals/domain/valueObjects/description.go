package valueobjects

type (
	Description struct {
		Description string
	}
)

func NewDescription(des string) (*Description, error) {

	return &Description{
		Description: des,
	}, nil
}
