package valueobjects

type (
	ProposalTimer struct {
		ProposalTimer int
	}
)

func NewProposalTimer(proposalTimer int) (ProposalTimer, error) {
	return ProposalTimer{
		ProposalTimer: proposalTimer,
	}, nil
}
