package errors

type roomNotFoundConst string

const ErrRoomNotFound roomNotFoundConst = "room not found."

func (r roomNotFoundConst) Error() string {
	return string(r)
}
