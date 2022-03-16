package whatsrhymen

import (
	"log"
	"time"

	whatsrhymen "github.com/Rhymen/go-whatsapp"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// Cria uma mensagem no formato do QuePasa apartir de uma msg do WhatsApp
// Preenche somente as propriedades padrões e comuns a todas as msgs
func CreateMessageFromInfo(Info whatsrhymen.MessageInfo) (message *whatsapp.WhatsappMessage) {
	message = &whatsapp.WhatsappMessage{}
	message.ID = Info.Id
	message.Timestamp = time.Unix(int64(Info.Timestamp), 0)

	// Fui eu quem enviou a msg ?
	message.FromMe = Info.FromMe
	return
}

func FillHeader(message *whatsapp.WhatsappMessage, Info whatsrhymen.MessageInfo, Conn *whatsrhymen.Conn) (err error) {

	// Endereço correto para onde deve ser devolvida a msg
	message.Chat.ID = Info.RemoteJid
	message.Chat.Title = getTitle(Conn.Store, Info.RemoteJid)

	// Pessoa que enviou a msg dentro de um grupo
	if Info.Source.Participant != nil {
		message.Participant.ID = *Info.Source.Participant
		message.Participant.Title = getTitle(Conn.Store, *Info.Source.Participant)
	}

	return
}

// Retorna algum titulo válido apartir de um jid
func getTitle(store *whatsrhymen.Store, jid string) string {
	var result string
	contact, ok := store.Contacts[jid]
	if ok {
		result = getContactTitle(contact)
	}
	return result
}

// Retorna algum titulo válido apartir de um contato do whatsapp
func getContactTitle(contact whatsrhymen.Contact) string {
	var result string
	result = contact.Name
	if len(result) == 0 {
		result = contact.Notify
		if len(result) == 0 {
			result = contact.Short
		}
	}
	return result
}

func FillImageAttachment(message *whatsapp.WhatsappMessage, msg whatsrhymen.ImageMessage, con *whatsrhymen.Conn) {
	if msg.Info.Source.Message.ImageMessage.Url == nil {
		// Aconteceu na primeira vez, quando cadastrei o número de whatsapp errado
		log.Println("erro on filling image attachement, url not avail")
		return
	}

	// getKey := msg.Info.Source.Message.ImageMessage.MediaKey
	getUrl := *msg.Info.Source.Message.ImageMessage.Url
	getLength := *msg.Info.Source.Message.ImageMessage.FileLength
	getMIME := *msg.Info.Source.Message.ImageMessage.Mimetype

	message.Attachment = &whatsapp.WhatsappAttachment{
		FileLength: getLength,
		Mimetype:   getMIME,
		FileName:   getUrl,
	}
}

func FillAudioAttachment(message *whatsapp.WhatsappMessage, msg whatsrhymen.AudioMessage, con *whatsrhymen.Conn) {
	// getKey := msg.Info.Source.Message.AudioMessage.MediaKey
	getUrl := *msg.Info.Source.Message.AudioMessage.Url
	getLength := *msg.Info.Source.Message.AudioMessage.FileLength
	getMIME := *msg.Info.Source.Message.AudioMessage.Mimetype

	message.Attachment = &whatsapp.WhatsappAttachment{
		FileLength: getLength,
		Mimetype:   getMIME,
		FileName:   getUrl,
	}
}

func FillDocumentAttachment(message *whatsapp.WhatsappMessage, msg whatsrhymen.DocumentMessage, con *whatsrhymen.Conn) {
	innerMSG := msg.Info.Source.Message.DocumentMessage

	message.Attachment = &whatsapp.WhatsappAttachment{
		FileLength: *innerMSG.FileLength,

		// Adicionando document no final do mime para saber que foi enviado como documento pelo whatsapp
		// Acontece de enviarem imagens como documento e não como imagens
		// Essa informação adicional é importante para realizar o download da media depois
		Mimetype: msg.Type + "; wa-document",
		FileName: msg.FileName,
	}
}

func Download(url string, mediakey []byte, media whatsrhymen.MediaType, length int) (data []byte, err error) {
	data, err = whatsrhymen.Download(url, mediakey, media, length)
	return
}
