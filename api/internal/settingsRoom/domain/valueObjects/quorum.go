package valueobjects

type (
	Quorum struct {
		Quorum *int //esta en
	}
)

func NewQuorum(quorum *int) (*Quorum, error) {

	return &Quorum{
		Quorum: quorum,
	}, nil
}
