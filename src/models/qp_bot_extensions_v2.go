package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
)

// returning []QPMessageV1
// bot.GetMessages(searchTime)
func GetMessagesFromBotV2(source QPBot, timestamp string) (messages []QPMessageV2, err error) {

	server, err := GetServerFromBot(source)
	if err != nil {
		return
	}

	searchTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		if len(timestamp) > 0 {
			return
		} else {
			searchTimestamp = 0
		}
	}

	searchTime := time.Unix(searchTimestamp, 0)
	messages = GetMessagesFromServerV2(server, searchTime)
	return
}

func ToQPBotV2(source *QPBot) (destination *QPBotV2) {
	destination = &QPBotV2{}
	err := copier.Copy(destination, source)
	if err != nil {
		log.Errorf("error on convert bot to version 1: %s", err.Error())
	}

	if !strings.Contains(destination.ID, "@") {
		destination.ID = destination.ID + "@c.us"
	}
	return
}
