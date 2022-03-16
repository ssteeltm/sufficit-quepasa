package models

import (
	log "github.com/sirupsen/logrus"

	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	whatsrhymen "github.com/sufficit/sufficit-quepasa-fork/whatsrhymen"
)

func NewWhatsrhymenConnection(wid string, logger *log.Logger) (whatsapp.IWhatsappConnection, error) {
	return whatsrhymen.WhatsrhymenService.CreateConnection(wid, logger)
}
