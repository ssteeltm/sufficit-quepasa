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
	} else if in.StickerMessage != nil {
		HandleStickerMessage(out, in.StickerMessage)
	} else if in.DocumentMessage != nil {
		HandleDocumentMessage(out, in.DocumentMessage)
	} else if in.AudioMessage != nil {
		HandleAudioMessage(out, in.AudioMessage)
	} else if in.VideoMessage != nil {
		HandleVideoMessage(out, in.VideoMessage)
	} else if in.ExtendedTextMessage != nil {
		HandleExtendedTextMessage(out, in.ExtendedTextMessage)
	} else if in.ProtocolMessage != nil || in.SenderKeyDistributionMessage != nil {
		out.Type = DiscardMessageType
	} else if len(out.Text) > 0 {
		out.Type = TextMessageType
	}
}

func HandleUnknownMessage(in interface{}) {
	log.Debug("Received an unknown message !")
	b, err := json.Marshal(in)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

}

// Msg em resposta a outra
func HandleExtendedTextMessage(out *WhatsappMessage, in *ExtendedTextMessage) {
	log.Debug("Received a text|extended message !")
	out.Type = TextMessageType

	if in.Text != nil {
		out.Text = *in.Text
	}

	info := in.ContextInfo
	if info != nil {
		if info.ForwardingScore != nil {
			out.ForwardingScore = *info.ForwardingScore
		}

		if info.StanzaId != nil {
			out.InReply = *info.StanzaId  
		}
	}
}

func HandleImageMessage(out *WhatsappMessage, in *ImageMessage) {
	log.Debug("Received an image message !")
	out.Content = in
	out.Type = ImageMessageType

	// in case of caption passed
	if in.Caption != nil {
		out.Text = *in.Caption
	}

	jpeg := GetStringFromBytes(in.JpegThumbnail)
	out.Attachment = &WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		JpegThumbnail: jpeg,
	}
}

func HandleStickerMessage(out *WhatsappMessage, in *StickerMessage) {
	log.Debug("Received a image|sticker message !")
	out.Content = in
	out.Type = ImageMessageType

	jpeg := GetStringFromBytes(in.PngThumbnail)
	out.Attachment = &WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		JpegThumbnail: jpeg,
	}
}

func HandleVideoMessage(out *WhatsappMessage, in *VideoMessage) {
	log.Debug("Received a video message !")
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
	log.Debug("Received a document message !")
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
