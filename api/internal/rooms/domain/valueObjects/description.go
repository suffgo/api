package valueobjects


type (
	Description struct {
		Description string
	}
)

func NewDescription(description string) (*Description, error) {


	return &Description{
		Description: description,
	}, nil
}