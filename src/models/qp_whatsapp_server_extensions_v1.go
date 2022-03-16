package models

import "time"

// returning []QPMessageV1
// server.GetMessages(searchTime)
func GetMessagesFromServerV1(server *QPWhatsappServer, searchTime time.Time) (messages []QPMessageV1, err error) {
	list, err := server.GetMessages(searchTime)
	if err != nil {
		return
	}

	for _, item := range list {
		messages = append(messages, ToQPMessageV1(item, server.GetWid()))
	}

	return
}
