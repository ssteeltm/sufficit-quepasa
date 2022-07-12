package models

import (
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

func SendMessageFromBot(source *QPBot, msg *whatsapp.WhatsappMessage) (err error) {
	server, err := GetServerFromBot(*source)
	if err != nil {
		return
	}

	return server.SendMessage(msg)
}
