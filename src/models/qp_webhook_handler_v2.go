package models

import (
	log "github.com/sirupsen/logrus"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QPWebhookHandlerV2 struct {
	Server *QPWhatsappServer
}

func (w QPWebhookHandlerV2) Handle(payload WhatsappMessage) {
	if payload.Type == UnknownMessageType {
		log.Debug("ignoring unknown message type on webhook request")
		return
	}

	msg := ToQPMessageV2(payload, w.Server.GetWid())
	PostToWebHookFromServer(w.Server, msg)
}
