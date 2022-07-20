package models

import (
	"mime"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type QpSendRequest struct {
	ChatId    string `json:"chatid"`
	FileName  string `json:"filename,omitempty"`
	TextLabel string `json:"textlabel,omitempty"`
	Content   []byte
}

func (source *QpSendRequest) ToWhatsappAttachment() (attach *whatsapp.WhatsappAttachment, err error) {
	attach = &whatsapp.WhatsappAttachment{}

	mimeType := http.DetectContentType(source.Content)
	log.Tracef("detected mime type: %s", mimeType)

	fileName := source.FileName
	// Defining a filename if not found before
	if len(fileName) == 0 {
		const layout = "20060201150405"
		t := time.Now().UTC()
		fileName = "file-" + t.Format(layout)

		// get file extension from mime type
		extension, _ := mime.ExtensionsByType(mimeType)
		if len(extension) > 0 {
			fileName = fileName + extension[0]
		}
	}

	attach.FileName = fileName
	attach.FileLength = uint64(len(source.Content))
	attach.Mimetype = mimeType
	attach.SetContent(&source.Content)
	return
}
