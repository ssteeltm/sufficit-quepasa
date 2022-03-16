package models

import (
	"strings"

	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

func SendMessageFromBot(source *QPBot, msg *WhatsappMessage) (err error) {
	server, err := GetServerFromBot(*source)
	if err != nil {
		return
	}

	return server.SendMessage(msg)
}

func ToQPBotV1(source *QPBot) (destination *QPBotV1) {
	destination = &QPBotV1{}
	err := copier.Copy(destination, source)
	if err != nil {
		log.Errorf("error on convert bot to version 1: %s", err.Error())
	}

	if !strings.Contains(destination.ID, "@") {
		destination.ID = destination.ID + "@c.us"
	}
	return
}
