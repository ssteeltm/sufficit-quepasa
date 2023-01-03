package whatsmeow

import (
	"fmt"
)

type WhatsmeowStoreNotFoundException struct {
	Wid string
}

func (e *WhatsmeowStoreNotFoundException) Error() string {
	return fmt.Sprintf("cant find a store for wid (%s)", e.Wid)
}

func (e *WhatsmeowStoreNotFoundException) Unauthorized() bool {
	return true
}
