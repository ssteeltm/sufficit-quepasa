package whatsapp

import (
	log "github.com/sirupsen/logrus"
)

type IWhatsappConnection interface {
	// Returns Connection Version (beta|multi|single)
	GetVersion() string

	GetStatus() WhatsappConnectionState

	// Retorna o ID do controlador whatsapp
	GetWid() (string, error)
	GetTitle() string

	Connect() error
	Disconnect() error
	Delete() error
	GetWhatsAppQRChannel(chan string) error
	UpdateHandler(IWhatsappHandlers)

	// Download message attachment if exists
	Download(IWhatsappMessage) (data []byte, err error)

	// Default send message method
	Send(WhatsappMessage) (IWhatsappSendResponse, error)

	// Define the log level for this connection
	LogLevel(log.Level)
	PrintStatus()
}
