package whatsmeow

import (
	"encoding/base64"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	. "go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

type WhatsmeowLogLevel string

const (
	DebugLevel WhatsmeowLogLevel = "DEBUG"
	InfoLevel  WhatsmeowLogLevel = "INFO"
	WarnLevel  WhatsmeowLogLevel = "WARN"
	ErrorLevel WhatsmeowLogLevel = "ERROR"
)

func GetMediaTypeFromAttachment(source *WhatsappAttachment) MediaType {
	return GetMediaTypeFromString(source.Mimetype)
}

// Traz o MediaType para download do whatsapp
func GetMediaTypeFromString(Mimetype string) MediaType {

	// usado pela API para garantir o envio como documento de qualquer anexo
	if strings.Contains(Mimetype, "wa-document") {
		return MediaDocument
	}

	// apaga informações após o ;
	// fica somente o mime mesmo
	mimeOnly := strings.Split(Mimetype, ";")
	switch mimeOnly[0] {
	case "image/png", "image/jpeg":
		return MediaImage
	case "audio/ogg", "audio/mpeg", "audio/mp4", "audio/x-wav":
		return MediaAudio
	case "video/mp4":
		return MediaVideo
	default:
		return MediaDocument
	}
}

func ToWhatsmeowMessage(source IWhatsappMessage) (msg *waProto.Message, err error) {
	messageText := source.GetText()

	if !source.HasAttachment() {
		internal := &waProto.ExtendedTextMessage{Text: &messageText}
		msg = &waProto.Message{ExtendedTextMessage: internal}
	}

	return
}

func NewWhatsmeowMessageAttachment(response UploadResponse, attach *WhatsappAttachment) (msg *waProto.Message) {
	media := GetMediaTypeFromString(attach.Mimetype)
	switch media {
	case MediaImage:
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
	case MediaAudio:
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
	case MediaVideo:
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
