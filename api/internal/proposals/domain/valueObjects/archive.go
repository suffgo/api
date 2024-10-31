package valueobjects

type (
	Archive struct {
		Archive string
	}
)

func NewArchive(archive string) (*Archive, error) {

	return &Archive{
		Archive: archive,
	}, nil
}
