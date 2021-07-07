package datatypes

import "fmt"

type NotFound struct {
	Message string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("not found: %s", e.Message)
}
