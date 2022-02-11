package controllers

import (
	"fmt"

	. "github.com/sufficit/sufficit-quepasa-fork/models"
)

type ApiServerNotReadyException struct {
	Wid    string
	Status QPWhatsappState
}

func (e *ApiServerNotReadyException) Error() string {
	return fmt.Sprintf("bot (%s) not ready yet ! current status: %s.", e.Wid, e.Status.String())
}
