package whatsrhymen

import (
	"mime"
	"time"

	whatsrhymen "github.com/Rhymen/go-whatsapp"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// Cria uma mensagem no formato do QuePasa apartir de uma msg do WhatsApp
// Preenche somente as propriedades padrões e comuns a todas as msgs
func CreateMessageFromInfo(Info whatsrhymen.MessageInfo) (message *WhatsrhymenMessage) {
	message = &WhatsrhymenMessage{}
	message.ID = Info.Id
	message.Timestamp = time.Unix(int64(Info.Timestamp), 0)

	// Fui eu quem enviou a msg ?
	message.FromMe = Info.FromMe
	return
}

func FillHeader(message *WhatsrhymenMessage, Info whatsrhymen.MessageInfo, Conn *whatsrhymen.Conn) (err error) {

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

func FillImageAttachment(message *WhatsrhymenMessage, msg whatsrhymen.ImageMessage, con *whatsrhymen.Conn) {
	innerMSG := msg.Info.Source.Message.ImageMessage

	getKey := innerMSG.GetMediaKey()
	getUrl := innerMSG.GetUrl()
	getLength := innerMSG.GetFileLength()
	getMIME := innerMSG.GetMimetype()
	getFileName := GetFileName(message.ID, getMIME, whatsrhymen.MediaImage)

	message.Attachment = &whatsapp.WhatsappAttachment{
		FileName:   getFileName,
		FileLength: getLength,
		Mimetype:   getMIME,
	}

	message.AttachmentInfo = &WhatsrhymenAttachmentInfo{
		Url:       getUrl,
		MediaKey:  getKey,
		Length:    int(getLength),
		MediaType: whatsrhymen.MediaImage,
	}
}

func FillAudioAttachment(message *WhatsrhymenMessage, msg whatsrhymen.AudioMessage, con *whatsrhymen.Conn) {
	innerMSG := msg.Info.Source.Message.AudioMessage

	getKey := innerMSG.GetMediaKey()
	getUrl := innerMSG.GetUrl()
	getLength := innerMSG.GetFileLength()
	getMIME := innerMSG.GetMimetype()
	getFileName := GetFileName(message.ID, getMIME, whatsrhymen.MediaAudio)

	message.Attachment = &whatsapp.WhatsappAttachment{
		FileName:   getFileName,
		FileLength: getLength,
		Mimetype:   getMIME,
	}

	message.AttachmentInfo = &WhatsrhymenAttachmentInfo{
		Url:       getUrl,
		MediaKey:  getKey,
		Length:    int(getLength),
		MediaType: whatsrhymen.MediaAudio,
	}
}

// using mime database to apply an extension on filename
func GetFileName(id string, mimetype string, mediatype whatsrhymen.MediaType) string {
	extensions, err := mime.ExtensionsByType(mimetype)
	if err != nil {
		if mediatype == whatsrhymen.MediaAudio {
			return id + ".ogg"
		} else if mediatype == whatsrhymen.MediaImage {
			return id + ".jpg"
		} else {
			return ""
		}
	} else {
		return id + extensions[0]
	}
}

func FillDocumentAttachment(message *WhatsrhymenMessage, msg whatsrhymen.DocumentMessage, con *whatsrhymen.Conn) {
	innerMSG := msg.Info.Source.Message.DocumentMessage

	getKey := innerMSG.GetMediaKey()
	getUrl := innerMSG.GetUrl()
	getLength := innerMSG.GetFileLength()
	getMIME := innerMSG.GetMimetype()

	message.Attachment = &whatsapp.WhatsappAttachment{
		FileLength: getLength,

		// Adicionando document no final do mime para saber que foi enviado como documento pelo whatsapp
		// Acontece de enviarem imagens como documento e não como imagens
		// Essa informação adicional é importante para realizar o download da media depois
		Mimetype: getMIME + "; wa-document",
		FileName: msg.FileName,
	}

	message.AttachmentInfo = &WhatsrhymenAttachmentInfo{
		Url:       getUrl,
		MediaKey:  getKey,
		Length:    int(getLength),
		MediaType: whatsrhymen.MediaDocument,
	}
}

func DownloadAttachment(url string, mediakey []byte, media whatsrhymen.MediaType, length int) (data []byte, err error) {
	return whatsrhymen.Download(url, mediakey, media, length)
}
