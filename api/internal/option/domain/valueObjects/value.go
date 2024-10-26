package valueobjects

type (
	Value struct{
		Value string
	}
)

func NewValue(value string) (*Value, error){

	return &Value{
		Value: value,
	},nil
}