package models

import (
	"strconv"
	"time"
)

// returning []QPMessageV1
// bot.GetMessages(searchTime)
func GetMessagesFromBotV1(source QPBot, timestamp string) (messages []QPMessageV1, err error) {

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
	return GetMessagesFromServerV1(server, searchTime)
}
