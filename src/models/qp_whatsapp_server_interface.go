package models

type IQPWhatsappServer interface {

	// Returns whatsapp controller id on E164
	GetWid() string

	// Download message attachments
	Download(id string) ([]byte, error)
}
