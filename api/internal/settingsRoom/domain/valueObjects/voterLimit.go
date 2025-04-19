package valueobjects

type (
	VoterLimit struct {
		VoterLimit int
	}
)

func NewVoterLimit(voterLimit int) (VoterLimit, error) {
	return VoterLimit{
		VoterLimit: voterLimit,
	}, nil
}
