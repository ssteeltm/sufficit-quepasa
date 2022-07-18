package whatsmeow

import (
	"encoding/base64"

	_ "github.com/mattn/go-sqlite3"

	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	whatsmeow "go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

type WhatsmeowLogLevel string

const (
	DebugLevel WhatsmeowLogLevel = "DEBUG"
	InfoLevel  WhatsmeowLogLevel = "INFO"
	WarnLevel  WhatsmeowLogLevel = "WARN"
	ErrorLevel WhatsmeowLogLevel = "ERROR"
)

func GetMediaTypeFromAttachment(source *whatsapp.WhatsappAttachment) whatsmeow.MediaType {
	return GetMediaTypeFromString(source.Mimetype)
}

// Traz o MediaType para download do whatsapp
func GetMediaTypeFromString(Mimetype string) whatsmeow.MediaType {

	msgType := whatsapp.GetMessageType(Mimetype)

	switch msgType {
	case whatsapp.ImageMessageType:
		return whatsmeow.MediaImage
	case whatsapp.AudioMessageType:
		return whatsmeow.MediaAudio
	case whatsapp.VideoMessageType:
		return whatsmeow.MediaVideo
	default:
		return whatsmeow.MediaDocument
	}
}

func ToWhatsmeowMessage(source whatsapp.IWhatsappMessage) (msg *waProto.Message, err error) {
	messageText := source.GetText()

	if !source.HasAttachment() {
		internal := &waProto.ExtendedTextMessage{Text: &messageText}
		msg = &waProto.Message{ExtendedTextMessage: internal}
	}

	return
}

func NewWhatsmeowMessageAttachment(response whatsmeow.UploadResponse, attach *whatsapp.WhatsappAttachment) (msg *waProto.Message) {
	media := GetMediaTypeFromString(attach.Mimetype)
	switch media {
	case whatsmeow.MediaImage:
		msg = &waProto.Message{ImageMessage: &waProto.ImageMessage{
			Url:           &response.URL,
			DirectPath:    &response.DirectPath,
			MediaKey:      response.MediaKey,
			FileEncSha256: response.FileEncSHA256,
			FileSha256:    response.FileSHA256,
			FileLength:    &response.FileLength,

			Mimetype: &attach.Mimetype,
			Caption:  &attach.FileName,
		},
		}
		return
	case whatsmeow.MediaAudio:
		internal := &waProto.AudioMessage{
			Url:           &response.URL,
			DirectPath:    &response.DirectPath,
			MediaKey:      response.MediaKey,
			FileEncSha256: response.FileEncSHA256,
			FileSha256:    response.FileSHA256,
			FileLength:    &response.FileLength,

			Mimetype: &attach.Mimetype,
			Ptt:      &[]bool{true}[0],
		}
		msg = &waProto.Message{AudioMessage: internal}
		return
	case whatsmeow.MediaVideo:
		internal := &waProto.VideoMessage{
			Url:           &response.URL,
			DirectPath:    &response.DirectPath,
			MediaKey:      response.MediaKey,
			FileEncSha256: response.FileEncSHA256,
			FileSha256:    response.FileSHA256,
			FileLength:    &response.FileLength,

			Mimetype: &attach.Mimetype,
			Caption:  &attach.FileName,
		}
		msg = &waProto.Message{VideoMessage: internal}
		return
	default:
		internal := &waProto.DocumentMessage{
			Url:           &response.URL,
			DirectPath:    &response.DirectPath,
			MediaKey:      response.MediaKey,
			FileEncSha256: response.FileEncSHA256,
			FileSha256:    response.FileSHA256,
			FileLength:    &response.FileLength,

			Mimetype: &attach.Mimetype,
			FileName: &attach.FileName,
		}
		msg = &waProto.Message{DocumentMessage: internal}
		return
	}
}

func GetStringFromBytes(bytes []byte) string {
	if bytes != nil {
		return base64.StdEncoding.EncodeToString(bytes)
	}
	return ""
}
