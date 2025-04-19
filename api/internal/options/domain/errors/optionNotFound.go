package errors

type OptionNotFoundErrorConst string

const ErrOptNotFound OptionNotFoundErrorConst = "option not found."

func (o OptionNotFoundErrorConst) Error() string {
	return string(o)
}
