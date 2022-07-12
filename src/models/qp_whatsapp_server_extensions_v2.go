package models

import "time"

// returning []QPMessageV1
// server.GetMessages(searchTime)
func GetMessagesFromServerV2(server *QPWhatsappServer, searchTime time.Time) (messages []QPMessageV2) {
	list := server.GetMessages(searchTime)
	for _, item := range list {
		messages = append(messages, ToQPMessageV2(item, server.GetWid()))
	}

	return
}
