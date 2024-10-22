package valueobjects


type (
	UserID struct{
		Id uint
	}
)


func NewUserID(id uint) *UserID {
	return &UserID{
		Id : id,
	}
}
