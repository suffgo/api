package errors

type voteNotFoundConst string

const ErrVoteNotFound voteNotFoundConst = "vote not found."

func (v voteNotFoundConst) Error() string {
	return string(v)
}
