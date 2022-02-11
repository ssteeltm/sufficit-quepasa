package models

import (
	"strings"

	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsmeow"
)

func NewConnection(wid string) (IWhatsappConnection, error) {
	return NewWhatsappConnection(wid)
}

func ToQPMessageV1(source WhatsappMessage, wid string) (message QPMessageV1) {
	message.ID = source.ID
	message.Timestamp = uint64(source.Timestamp.Unix())
	message.Text = source.Text
	message.FromMe = source.FromMe

	WhatsappID := wid
	if !strings.Contains(WhatsappID, "@") {
		WhatsappID = WhatsappID + "@s.whatsapp.net"
	}

	message.Controller = QPEndPoint{ID: WhatsappID}
	message.ReplyTo = ChatToQPEndPoint(source.Chat)

	if (WhatsappEndpoint{}) != source.Participant {
		message.Participant = ToQPEndPoint(source.Participant)
	}

	if source.HasAttachment() {
		message.Attachment = ToQPAttachment(source.Attachment, message.ID, WhatsappID)
	}

	return
}
