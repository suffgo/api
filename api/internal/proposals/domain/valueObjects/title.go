package valueobjects

type (
	Title struct {
		Title string
	}
)

func NewTitle(title string) (*Title, error) {

	return &Title{
		Title: title,
	}, nil
}
