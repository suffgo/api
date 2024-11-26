package valueobjects

type (
	Quorum struct {
		Quorum *int
	}
)

func NewQuorum(quorum *int) (*Quorum, error) {

	return &Quorum{
		Quorum: quorum,
	}, nil
}
