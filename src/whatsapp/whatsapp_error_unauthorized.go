package whatsapp

import (
	"fmt"
)

type UnauthorizedError struct {
	Inner error
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("Unauthorized: %s", e.Inner)
}
