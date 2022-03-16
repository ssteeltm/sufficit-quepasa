package whatsapp

type IWhatsappMessage interface {

	// E164 Phone without trailing + or GroupID with -
	// Ex: 5521967609095
	// Ex: 5521967609095-1445779956
	GetChatID() string

	// Clear text message or html encoded
	GetText() string

	// Check if that msg has a valid attachment
	HasAttachment() bool

	// Original message from source service
	GetSource() interface{}
}
