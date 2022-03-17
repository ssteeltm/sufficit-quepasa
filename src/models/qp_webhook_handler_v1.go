package models

import (
	log "github.com/sirupsen/logrus"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QPWebhookHandlerV1 struct {
	Server *QPWhatsappServer
}

func (w *QPWebhookHandlerV1) Handle(payload WhatsappMessage) {
	if !w.HasWebhook() {
		return
	}

	if payload.Type == UnknownMessageType {
		log.Debug("ignoring unknown message type on webhook request")
		return
	}

	if payload.Type == TextMessageType && len(payload.Text) == 0 {
		log.Warnf("ignoring empty text message on webhook request: %v", payload.ID)
		return
	}

	msg := ToQPMessageV1(payload, w.Server.GetWid())
	PostToWebHookFromServer(w.Server, msg)
}

func (w *QPWebhookHandlerV1) HasWebhook() bool {
	if w.Server != nil {
		url := w.Server.Webhook()
		if len(url) > 0 {
			return true
		}
	}
	return false
}
