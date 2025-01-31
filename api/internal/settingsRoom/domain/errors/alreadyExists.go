package errors

type alreadyExists string

const ErrAlreadyExists alreadyExists = "extended settings room already exists."

func (r alreadyExists) Error() string {
	return string(r)
}