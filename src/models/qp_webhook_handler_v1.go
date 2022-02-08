package models

import (
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QPWebhookHandlerV1 struct {
	Server *QPWhatsappServer
}

func (w QPWebhookHandlerV1) Handle(payload WhatsappMessage) {
	msg := ToQPMessageV1(payload, w.Server.GetWid())
	PostToWebHookFromServer(w.Server, msg)
}
