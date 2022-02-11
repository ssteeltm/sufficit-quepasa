package models

import (
	"encoding/base64"

	. "github.com/sufficit/sufficit-quepasa-fork/library"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

func ToWhatsappAttachment(source *QPAttachmentV1) (attach *WhatsappAttachment, err error) {
	attach = &WhatsappAttachment{}
	content, err := base64.StdEncoding.DecodeString(source.Base64)
	if err != nil {
		return
	}

	attach.SetContent(&content)
	attach.FileName = source.FileName
	attach.Mimetype = source.MIME
	attach.FileLength = uint64(source.Length)
	return
}

func ToWhatsappMessageV1(source *QPSendRequestV1) (msg *WhatsappMessage, err error) {
	recipient, err := FormatEndpoint(source.Recipient)
	if err != nil {
		return
	}

	attach, err := ToWhatsappAttachment(&source.Attachment)
	if err != nil {
		return
	}

	chat := WhatsappChat{ID: recipient}
	msg = &WhatsappMessage{}
	msg.Text = source.Message
	msg.Chat = chat
	if attach != nil {
		msg.Attachment = attach
	}
	return
}
