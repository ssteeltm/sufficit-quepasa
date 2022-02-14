package whatsmeow

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	. "go.mau.fi/whatsmeow/binary/proto"
)

func HandleKnowingMessages(out *WhatsappMessage, in *Message) {
	if in.ImageMessage != nil {
		HandleImageMessage(out, in.ImageMessage)
	} else if in.DocumentMessage != nil {
		HandleDocumentMessage(out, in.DocumentMessage)
	} else if in.AudioMessage != nil {
		HandleAudioMessage(out, in.AudioMessage)
	} else if in.VideoMessage != nil {
		HandleVideoMessage(out, in.VideoMessage)
	} else if in.ExtendedTextMessage != nil {
		HandleTextMessage(out, in.ExtendedTextMessage)
	} else if len(out.Text) > 0 {
		out.Type = TextMessageType
	}
}

func HandleTextMessage(out *WhatsappMessage, in *ExtendedTextMessage) {
	b, err := json.Marshal(in)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("LOGGING EXTENDED ::: " + string(b))
}

func HandleUnknownMessage(in interface{}) {

	b, err := json.Marshal(in)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

}

func HandleImageMessage(out *WhatsappMessage, in *ImageMessage) {
	log.Debug("Received an image message !")
	out.Content = in
	out.Type = ImageMessageType

	// in case of caption passed
	if in.Caption != nil {
		out.Text = *in.Caption
	}

	jpeg := base64.StdEncoding.EncodeToString(in.JpegThumbnail)
	out.Attachment = &WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		JpegThumbnail: jpeg,
	}
}

func HandleVideoMessage(out *WhatsappMessage, in *VideoMessage) {
	log.Debug("Received an video message !")
	out.Content = in
	out.Type = VideoMessageType

	// in case of caption passed
	if in.Caption != nil {
		out.Text = *in.Caption
	}

	jpeg := base64.StdEncoding.EncodeToString(in.JpegThumbnail)
	out.Attachment = &WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		JpegThumbnail: jpeg,
	}
}

func HandleDocumentMessage(out *WhatsappMessage, in *DocumentMessage) {
	log.Debug("Received an document message !")
	out.Content = in
	out.Type = DocumentMessageType

	out.Text = *in.Title

	jpeg := base64.StdEncoding.EncodeToString(in.JpegThumbnail)
	out.Attachment = &WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		FileName:      *in.FileName,
		JpegThumbnail: jpeg,
	}
}

func HandleAudioMessage(out *WhatsappMessage, in *AudioMessage) {
	log.Debug("Received an audio message !")
	out.Content = in
	out.Type = AudioMessageType

	out.Attachment = &WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		Seconds: *in.Seconds,
	}
}
