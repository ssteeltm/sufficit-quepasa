package models

import (
	"strings"

	log "github.com/sirupsen/logrus"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QPWebhookHandlerV2 struct {
	Server *QPWhatsappServer
}

func (w *QPWebhookHandlerV2) Handle(payload whatsapp.WhatsappMessage) {
	if !w.HasWebhook() {
		return
	}

	if payload.Type == 0 {
		log.Debug("ignoring unknown message type on webhook request")
		return
	}

	if payload.Type == whatsapp.TextMessageType && len(strings.TrimSpace(payload.Text)) <= 0 {
		log.Warnf("ignoring empty text message on webhook request: %v", payload.ID)
		return
	} else {
		log.Infof("message on webhook request: %v, type: %v", payload.ID, payload.Type)
	}

	if payload.Chat.ID == "status@broadcast" && !w.Server.HandleBroadcast() {
		log.Warnf("ignoring broadcast message on webhook request: %v", payload.ID)
		return
	}

	msg := ToQPMessageV2(payload, w.Server.GetWid())
	PostToWebHookFromServer(w.Server, msg)
}

func (w *QPWebhookHandlerV2) HasWebhook() bool {
	if w.Server != nil {
		url := w.Server.Webhook()
		if len(url) > 0 {
			return true
		}
	}
	return false
}
