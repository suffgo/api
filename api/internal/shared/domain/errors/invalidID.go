package errors

import "fmt"

type InvalidIDError struct {
	ID string
}

func (e *InvalidIDError) Error() string {
	return fmt.Sprintf("Invalid id %s", e.ID)
}