package errors

type RepeatedOptionErrorConst string

const ErrOptRepeated RepeatedOptionErrorConst = "option value already taken."

func (o RepeatedOptionErrorConst) Error() string {
	return string(o)
}
