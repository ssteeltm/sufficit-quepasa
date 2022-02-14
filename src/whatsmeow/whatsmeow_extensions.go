package whatsmeow

import (
	"encoding/base64"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"

	. "go.mau.fi/whatsmeow"
	//. "go.mau.fi/whatsmeow/types"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

type WhatsmeowLogLevel string

const (
	DebugLevel WhatsmeowLogLevel = "DEBUG"
	InfoLevel                    = "INFO"
	WarnLevel                    = "WARN"
	ErrorLevel                   = "ERROR"
)

// Flush entire Whatsmeow Database
// Use with wisdom !
func FlushDatabase() (err error) {
	devices, err := WhatsmeowService.Container.GetAllDevices()
	if err != nil {
		return
	}

	for _, element := range devices {
		err = element.Delete()
		if err != nil {
			return
		}
	}

	return
}

func NewWhatsappConnection(wid string, logger *log.Logger) (conn *WhatsmeowConnection, err error) {
	client, err := GetWhatsAppClient(wid, logger)
	if err != nil {
		return
	}

	handlers := &WhatsmeowHandlers{Client: client}
	err = handlers.Register()
	if err != nil {
		return
	}

	loggerEntry := logger.WithField("wid", wid)
	conn = &WhatsmeowConnection{Client: client, Handlers: handlers, Logger: logger, log: loggerEntry}
	return
}

func GetWhatsAppClient(wid string, logger *log.Logger) (client *Client, err error) {
	deviceStore, err := GetStoreFromWid(wid)
	if err != nil {
		err = fmt.Errorf("error on getting whatsapp client: %s", err)
	} else {
		clientLog := waLog.Stdout("Whatsmeow/Client", logger.Level.String(), true)
		client = NewClient(deviceStore, clientLog)
	}
	return
}

func GetStoreFromWid(wid string) (str *store.Device, err error) {
	if wid == "" {
		str = WhatsmeowService.Container.NewDevice()
	} else {
		devices, err := WhatsmeowService.Container.GetAllDevices()
		if err != nil {
			err = fmt.Errorf("error on getting store from wid (%s): %v", wid, err)
			return str, err
		} else {
			for _, element := range devices {
				if element.ID.User == wid {
					str = element
					break
				}
			}

			if str == nil {
				err = &WhatsmeowStoreNotFoundException{Wid: wid}
				return str, err
			}
		}
	}

	return
}

func FormatEndpoint(destination string) string {
	if strings.Contains(destination, "-") {
		return destination + "@g.us"
	} else {
		return destination + "@s.whatsapp.net"
	}
}

func GetMediaTypeFromAttachment(source *WhatsappAttachment) MediaType {
	return GetMediaType(source.Mimetype)
}

// Traz o MediaType para download do whatsapp
func GetMediaType(Mimetype string) MediaType {

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
	media := GetMediaType(attach.Mimetype)
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
