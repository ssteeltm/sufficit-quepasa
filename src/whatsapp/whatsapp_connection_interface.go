package whatsapp

type IWhatsappConnection interface {
	// Returns Connection Version (beta|multi|single)
	GetVersion() string

	// Retorna o ID do controlador whatsapp
	GetWid() (string, error)
	GetTitle() string

	Connect() error
	Disconnect() error
	GetWhatsAppQRChannel(chan string) error
	UpdateHandler(IWhatsappHandlers)

	// Download message attachment if exists
	Download(IWhatsappMessage) (data []byte, err error)

	// Default send message method
	Send(WhatsappMessage) (IWhatsappSendResponse, error)
}
