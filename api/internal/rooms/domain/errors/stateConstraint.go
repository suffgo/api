package errors

type roomStateConstraint string

const ErrStateConstraint roomStateConstraint = "the current room state doesnt support this operation."

func (r roomStateConstraint) Error() string {
	return string(r)
}