package models

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	. "github.com/sufficit/sufficit-quepasa-fork/library"
)

func SendMessageFromBOT(botID string, recipient string, text string, attachment QPAttachment) (messageID string, err error) {
	recipient, err = FormatEndpoint(recipient)
	if err != nil {
		return
	}

	server, err := GetServerFromID(botID)
	if err != nil {
		return
	}

	sendResponse, err := server.Send(recipient, text)
	return sendResponse.GetID(), err
}

// Retrieve messages from the controller, external
func RetrieveMessages(botID string, timestamp string) (messages []QPMessageV1, err error) {
	searchTime := time.Now().Add(-24 * time.Hour)

	searchTimestamp, _ := strconv.ParseInt(timestamp, 10, 64)
	if searchTimestamp != 0 {
		searchTime = time.Unix(searchTimestamp, 0)
	}

	server, ok := WhatsAppService.Servers[botID]
	if !ok {
		err = fmt.Errorf("handlers not read yet, please wait")
		return
	}

	messages, err = GetMessagesFromServerV1(server, searchTime)
	if err != nil {
		err = fmt.Errorf("msgs not read yet, please wait")
		return
	}

	sort.Sort(ByTimestampV1(messages))
	return
}
