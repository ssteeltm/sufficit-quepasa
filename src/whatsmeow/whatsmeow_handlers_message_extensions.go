package whatsmeow

import (
	"encoding/base64"
	"encoding/json"
	"mime"

	log "github.com/sirupsen/logrus"
	whatsapp "github.com/sufficit/sufficit-quepasa/whatsapp"
	. "go.mau.fi/whatsmeow/binary/proto"
)

func HandleKnowingMessages(handler *WhatsmeowHandlers, out *whatsapp.WhatsappMessage, in *Message) {
	if in.ImageMessage != nil {
		HandleImageMessage(handler.log, out, in.ImageMessage)
	} else if in.StickerMessage != nil {
		HandleStickerMessage(handler.log, out, in.StickerMessage)
	} else if in.DocumentMessage != nil {
		HandleDocumentMessage(handler.log, out, in.DocumentMessage)
	} else if in.AudioMessage != nil {
		HandleAudioMessage(handler.log, out, in.AudioMessage)
	} else if in.VideoMessage != nil {
		HandleVideoMessage(handler.log, out, in.VideoMessage)
	} else if in.ExtendedTextMessage != nil {
		HandleExtendedTextMessage(handler.log, out, in.ExtendedTextMessage)
	} else if in.ProtocolMessage != nil || in.SenderKeyDistributionMessage != nil {
		out.Type = whatsapp.DiscardMessageType
	} else if len(in.GetConversation()) > 0 {
		HandleTextMessage(handler.log, out, in)
	}
}

func HandleUnknownMessage(log *log.Entry, in interface{}) {
	log.Info("Received an unknown message !")
	b, err := json.Marshal(in)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug(string(b))
}

func HandleTextMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *Message) {
	log.Debug("Received a text message !")
	out.Type = whatsapp.TextMessageType
	out.Text = in.GetConversation()
}

// Msg em resposta a outra
func HandleExtendedTextMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *ExtendedTextMessage) {
	log.Debug("Received a text|extended message !")
	out.Type = whatsapp.TextMessageType

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

func HandleImageMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *ImageMessage) {
	log.Debug("Received an image message !")
	out.Content = in
	out.Type = whatsapp.ImageMessageType

	// in case of caption passed
	if in.Caption != nil {
		out.Text = *in.Caption
	}

	jpeg := GetStringFromBytes(in.JpegThumbnail)
	out.Attachment = &whatsapp.WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		JpegThumbnail: jpeg,
	}
}

func HandleStickerMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *StickerMessage) {
	log.Debug("Received a image|sticker message !")
	out.Content = in
	out.Type = whatsapp.ImageMessageType

	jpeg := GetStringFromBytes(in.PngThumbnail)
	out.Attachment = &whatsapp.WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		JpegThumbnail: jpeg,
	}
}

func HandleVideoMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *VideoMessage) {
	log.Debug("Received a video message !")
	out.Content = in
	out.Type = whatsapp.VideoMessageType

	// in case of caption passed
	if in.Caption != nil {
		out.Text = *in.Caption
	}

	jpeg := base64.StdEncoding.EncodeToString(in.JpegThumbnail)
	out.Attachment = &whatsapp.WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		JpegThumbnail: jpeg,
	}
}

func HandleDocumentMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *DocumentMessage) {
	log.Debug("Received a document message !")
	out.Content = in
	out.Type = whatsapp.DocumentMessageType

	if in.Title != nil {
		out.Text = *in.Title
	}

	jpeg := base64.StdEncoding.EncodeToString(in.JpegThumbnail)
	out.Attachment = &whatsapp.WhatsappAttachment{
		Mimetype:   *in.Mimetype + "; wa-document",
		FileLength: *in.FileLength,

		FileName:      *in.FileName,
		JpegThumbnail: jpeg,
	}
}

func HandleAudioMessage(log *log.Entry, out *whatsapp.WhatsappMessage, in *AudioMessage) {
	log.Debug("Received an audio message !")
	out.Content = in
	out.Type = whatsapp.AudioMessageType

	var seconds uint32
	if in.Seconds != nil {
		seconds = *in.Seconds
	}

	out.Attachment = &whatsapp.WhatsappAttachment{
		Mimetype:   *in.Mimetype,
		FileLength: *in.FileLength,

		Seconds: seconds,
	}

	// get file extension from mime type
	extension, _ := mime.ExtensionsByType(out.Attachment.Mimetype)
	if len(extension) > 0 {
		out.Attachment.FileName = out.Id + extension[0]
	}
}
