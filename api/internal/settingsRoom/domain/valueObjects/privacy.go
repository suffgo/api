package valueobjects

type (
	Privacy struct {
		Privacy bool
	}
)

func NewPrivacy(privacy bool) (Privacy, error) {

	return Privacy{
		Privacy: privacy,
	}, nil
}
