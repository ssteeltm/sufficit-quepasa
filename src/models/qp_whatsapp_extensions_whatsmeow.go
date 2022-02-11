package models

import (
	"strings"

	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsmeow"
)

func NewConnection(wid string) (IWhatsappConnection, error) {
	return NewWhatsappConnection(wid)
}

func ToQPMessageV2(source WhatsappMessage, wid string) (message QPMessageV2) {
	message.ID = source.ID
	message.Timestamp = uint64(source.Timestamp.Unix())
	message.Text = source.Text
	message.FromMe = source.FromMe

	if !strings.Contains(wid, "@") {
		message.Controller = QPEndpointV2{ID: wid + "@c.us"}
	} else {
		message.Controller = QPEndpointV2{ID: wid}
	}

	message.ReplyTo = ChatToQPEndPointV2(source.Chat)
	message.Chat = ChatToQPChatV2(source.Chat)

	if (WhatsappEndpoint{}) != source.Participant {
		message.Participant = ToQPEndPointV2(source.Participant)
	}

	if source.HasAttachment() {
		message.Attachment = ToQPAttachment(source.Attachment, message.ID, wid)
	}

	return
}

func ToQPMessageV1(source WhatsappMessage, wid string) (message QPMessageV1) {
	message.ID = source.ID
	message.Timestamp = uint64(source.Timestamp.Unix())
	message.Text = source.Text
	message.FromMe = source.FromMe

	if !strings.Contains(wid, "@") {
		message.Controller = QPEndpointV1{ID: wid + "@c.us"}
	} else {
		message.Controller = QPEndpointV1{ID: wid}
	}

	message.ReplyTo = ChatToQPEndPointV1(source.Chat)

	if (WhatsappEndpoint{}) != source.Participant {
		message.Participant = ToQPEndPointV1(source.Participant)
	}

	if source.HasAttachment() {
		message.Attachment = ToQPAttachment(source.Attachment, message.ID, wid)
	}

	return
}
