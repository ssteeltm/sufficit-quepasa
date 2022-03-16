package whatsrhymen

import (
	"bytes"
	"fmt"
	"strings"

	whatsrhymen "github.com/Rhymen/go-whatsapp"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// Default SEND method using WhatsappMessage Interface
func (conn *WhatsrhymenConnection) SendText(msg whatsapp.WhatsappMessage) (whatsapp.IWhatsappSendResponse, error) {
	response := &whatsapp.WhatsappSendResponse{}
	var err error
	messageText := msg.GetText()

	// Informações basicas para todo tipo de mensagens
	info := whatsrhymen.MessageInfo{
		RemoteJid: msg.Chat.ID,
	}

	if len(messageText) > 0 {
		wamsg := whatsrhymen.TextMessage{
			Info: info,
			Text: messageText,
		}
		response.ID, err = conn.Client.Send(wamsg)
	} else {
		err = fmt.Errorf("invalid text length")
	}

	return response, err
}

// Default SEND method using WhatsappMessage Interface
func (conn *WhatsrhymenConnection) SendAudio(msg whatsapp.WhatsappMessage) (whatsapp.IWhatsappSendResponse, error) {
	response := &whatsapp.WhatsappSendResponse{}
	var err error

	// Informações basicas para todo tipo de mensagens
	info := whatsrhymen.MessageInfo{
		RemoteJid: msg.Chat.ID,
	}

	content := *msg.Attachment.GetContent()
	if content != nil {

		// Definindo leitor de bytes do arquivo
		// Futuramente fazer download de uma URL para permitir tamanhos maiores
		reader := bytes.NewReader(content)

		ptt := strings.HasPrefix(msg.Attachment.Mimetype, "audio/ogg")
		wamsg := whatsrhymen.AudioMessage{
			Info:    info,
			Length:  uint32(msg.Attachment.FileLength),
			Type:    msg.Attachment.Mimetype,
			Ptt:     ptt,
			Content: reader,
		}
		response.ID, err = conn.Client.Send(wamsg)

	} else {
		err = fmt.Errorf("invalid content length")
	}

	return response, err
}

// Default SEND method using WhatsappMessage Interface
func (conn *WhatsrhymenConnection) SendImage(msg whatsapp.WhatsappMessage) (whatsapp.IWhatsappSendResponse, error) {
	response := &whatsapp.WhatsappSendResponse{}
	var err error
	messageText := msg.GetText()

	// Informações basicas para todo tipo de mensagens
	info := whatsrhymen.MessageInfo{
		RemoteJid: msg.Chat.ID,
	}

	content := *msg.Attachment.GetContent()
	if content != nil {
		// Definindo leitor de bytes do arquivo
		// Futuramente fazer download de uma URL para permitir tamanhos maiores
		reader := bytes.NewReader(content)
		wamsg := whatsrhymen.ImageMessage{
			Info:    info,
			Caption: messageText,
			Type:    msg.Attachment.Mimetype,
			Content: reader,
		}
		response.ID, err = conn.Client.Send(wamsg)

	} else {
		err = fmt.Errorf("invalid content length")
	}

	return response, err
}

// Default SEND method using WhatsappMessage Interface
func (conn *WhatsrhymenConnection) SendDocument(msg whatsapp.WhatsappMessage) (whatsapp.IWhatsappSendResponse, error) {
	response := &whatsapp.WhatsappSendResponse{}
	var err error

	// Informações basicas para todo tipo de mensagens
	info := whatsrhymen.MessageInfo{
		RemoteJid: msg.Chat.ID,
	}

	content := *msg.Attachment.GetContent()
	if content != nil {

		// Definindo leitor de bytes do arquivo
		// Futuramente fazer download de uma URL para permitir tamanhos maiores
		reader := bytes.NewReader(content)

		filename := msg.Attachment.FileName
		conn.log.Tracef("sending file, filename: %s", filename)

		wamsg := whatsrhymen.DocumentMessage{
			Info:    info,
			Title:   filename,
			Type:    msg.Attachment.Mimetype,
			Content: reader,
		}
		response.ID, err = conn.Client.Send(wamsg)
	} else {
		err = fmt.Errorf("invalid content length")
	}

	return response, err
}

// Traz o MediaType para download do whatsapp
func WAMediaType(m whatsapp.WhatsappAttachment) whatsrhymen.MediaType {

	if strings.Contains(m.Mimetype, "wa-document") {
		return whatsrhymen.MediaDocument
	}

	// apaga informações após o ;
	// fica somente o mime mesmo
	mimeOnly := strings.Split(m.Mimetype, ";")
	switch mimeOnly[0] {
	case "image/png", "image/jpeg":
		return whatsrhymen.MediaImage
	case "audio/ogg", "audio/mpeg", "audio/mp4", "audio/x-wav":
		return whatsrhymen.MediaAudio
	default:
		return whatsrhymen.MediaDocument
	}
}

func MessageTypeFromAttachment(m whatsapp.WhatsappAttachment) whatsapp.WhatsappMessageType {
	mediatype := WAMediaType(m)
	switch mediatype {
	case whatsrhymen.MediaImage:
		return whatsapp.ImageMessageType
	case whatsrhymen.MediaAudio:
		return whatsapp.AudioMessageType
	case whatsrhymen.MediaDocument:
		return whatsapp.DocumentMessageType
	default:
		return whatsapp.UnknownMessageType
	}
}
