package errors

type proposalNotFoundConst string

const ErrPropNotFound proposalNotFoundConst = "proposal not found."

func (p proposalNotFoundConst) Error() string {
	return string(p)
}
