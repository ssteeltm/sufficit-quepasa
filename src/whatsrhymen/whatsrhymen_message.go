package whatsrhymen

import (
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type WhatsrhymenMessage struct {
	whatsapp.WhatsappMessage
	AttachmentInfo *WhatsrhymenAttachmentInfo
}
