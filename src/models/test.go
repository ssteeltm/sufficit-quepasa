package models

import (
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
