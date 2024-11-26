package errors

type invalidIdConst string

const ErrInvalidID invalidIdConst = "invalid id."

func (u invalidIdConst) Error() string{
	return string(u)
}