package models

import (
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QPWebhookHandlerV2 struct {
	Server *QPWhatsappServer
}

func (w QPWebhookHandlerV2) Handle(payload WhatsappMessage) {
	msg := ToQPMessageV2(payload, w.Server.GetWid())
	PostToWebHookFromServer(w.Server, msg)
}
