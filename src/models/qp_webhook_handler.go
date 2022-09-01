package models

import (
	"strings"

	log "github.com/sirupsen/logrus"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
)

type QPWebhookHandler struct {
	Server *QPWhatsappServer
}

func (w *QPWebhookHandler) Handle(payload *whatsapp.WhatsappMessage) {
	if !w.HasWebhook() {
		return
	}

	if payload.Type == whatsapp.DiscardMessageType|whatsapp.UnknownMessageType {
		log.Debug("ignoring unknown message type on webhook request")
		return
	}

	if payload.Type == whatsapp.TextMessageType && len(strings.TrimSpace(payload.Text)) <= 0 {
		log.Debug("ignoring empty text message on webhook request: %v", payload.ID)
		return
	}

	if payload.Chat.ID == "status@broadcast" && !w.Server.HandleBroadcast() {
		log.Debug("ignoring broadcast message on webhook request: %v", payload.ID)
		return
	}

	PostToWebHookFromServer(w.Server, payload)
}

func (w *QPWebhookHandler) HasWebhook() bool {
	if w.Server != nil {
		return len(w.Server.Webhooks) > 0
	}
	return false
}
