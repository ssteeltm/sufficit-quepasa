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

		if len(messageText) == 0 {
			messageText := msg.Attachment.FileName
			if idx := strings.IndexByte(messageText, '.'); idx >= 0 {
				messageText = messageText[:idx]
			}
		}

		wamsg := whatsrhymen.DocumentMessage{
			Info:     info,
			Title:    messageText,
			FileName: msg.Attachment.FileName,
			Type:     msg.Attachment.Mimetype,
			Content:  reader,
		}
		response.ID, err = conn.Client.Send(wamsg)
	} else {
		err = fmt.Errorf("invalid content length")
	}

	return response, err
}
