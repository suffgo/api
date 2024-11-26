package errors

import "fmt"

//Este tipo de error se devuelve cuando hay un error de mapeo entre modelos y entidades
type DataMappingErrorConst string

const ErrDataMap DataMappingErrorConst = "data mapping error."

func (d DataMappingErrorConst) Error() string {
	return fmt.Sprintf("data mapping error")
}