package errors

import "fmt"

type InvalidValueError struct {
	Value string
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("el valor '%s' no es válido", e.Value)
}
