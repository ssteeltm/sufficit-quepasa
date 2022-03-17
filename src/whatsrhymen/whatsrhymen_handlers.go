package whatsrhymen

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	whatsrhymen "github.com/Rhymen/go-whatsapp"
	whatsapp "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type WhatsrhymenHandlers struct {
	Connection               *WhatsrhymenConnection
	unregisterRequestedToken bool
	attached                 bool
	log                      *log.Entry
}

func (handler *WhatsrhymenHandlers) Client() *whatsrhymen.Conn {
	return handler.Connection.Client
}

func (handler *WhatsrhymenHandlers) Store() *whatsrhymen.Store {
	return handler.Connection.Client.Store
}

// Aqui se processar um evento de recebimento de uma mensagem genérica
func (handler *WhatsrhymenHandlers) Message(message *WhatsrhymenMessage) {
	if handler.Connection.WAHandlers != nil {
		wamsg := &message.WhatsappMessage
		wamsg.Content = message

		// following to internal handlers
		go handler.Connection.WAHandlers.Message(wamsg)
	}
}

func (h *WhatsrhymenHandlers) IsRegistered() bool {
	return h.attached
}

func (h *WhatsrhymenHandlers) Register() (err error) {
	connection := h.Client()
	if connection == nil {
		return fmt.Errorf("not connected, on trying to register handlers")
	}

	if connection.Store == nil {
		err = fmt.Errorf("this client lost the store, probably a logout from whatsapp phone")
		return
	}

	h.unregisterRequestedToken = false
	connection.AddHandler(h)
	h.attached = true
	return
}

func (h *WhatsrhymenHandlers) UnRegister() {
	h.unregisterRequestedToken = true
	if h.Client().RemoveHandler(h) {
		h.attached = false
	}
}

// Essencial
// Unico item realmente necessario para o sistema do whatsapp funcionar
// Trata qualquer erro que influêncie no recebimento de msgs
func (h *WhatsrhymenHandlers) HandleError(publicError error) {
	if e, ok := publicError.(*whatsrhymen.ErrConnectionFailed); ok {
		// Erros comuns de desconexão por qualquer motivo aleatório
		if strings.Contains(e.Err.Error(), "close 1006") {
			// 1006 falha no websocket, informações inválidas, provavelmente baixa qualidade de internet no celular
			h.log.Printf("Websocket corrupted (probably poor connection quality), should restart ...")
		} else if strings.Contains(e.Err.Error(), "connection reset by peer") {
			// Possível falha no keep alive, muito tempo sem tráfego
			h.log.Printf("Websocket corrupted (probably inactive), should restart ...")
		} else {
			h.log.Error(e.Err)
		}
		return
	} else if strings.Contains(publicError.Error(), "code: 1000") {

		// Desconexão forçado é algum evento iniciado pelo whatsapp
		// tested with user logout device on whatsapp
		// tested with mobile phone lost internet connection
		h.log.Errorf("forced by whatsapp: %s", publicError.Error())
		h.Connection.Failed("handlers")
		return
	} else if strings.Contains(publicError.Error(), "close 1006") {
		// Desconexão forçado é algum evento iniciado pelo whatsapp
		h.log.Printf("Desconexão por falha no websocket, code: 1006, iremos reiniciar automaticamente")
		// Se houve desconexão, reseta
		return
	} else if strings.Contains(publicError.Error(), "keepAlive failed") {
		// Se houve desconexão, reseta
		h.log.Printf("Keep alive failed, waiting ...")
		// go h.Server.Restart()
		return
	} else if strings.Contains(publicError.Error(), "server closed connection") {
		// Se houve desconexão, reseta
		h.log.Printf("Server closed connection, restarting ...")
		return
	} else if strings.Contains(publicError.Error(), "message type not implemented") {
		// Ignorando, nova implementação com Handlers não criados ainda
		return
	} else {
		h.log.Error(publicError)
	}
}

func (h *WhatsrhymenHandlers) HandleJsonMessage(message string) {
	var waJsonMessage WhatsrhymenMessageJson
	err := json.Unmarshal([]byte(message), &waJsonMessage)
	if err == nil {
		if waJsonMessage.Cmd.Type == "disconnect" {
			// Restarting because an order of whatsapp
			h.log.Printf("Restart Order by: %s", waJsonMessage.Cmd.Kind)
		} else {
			h.log.Debug("(DEV) JSON Unmarshal string :: %s", message)
			h.log.Debug("(DEV) JSON Unmarshal :: %s", waJsonMessage)
		}
	} else {
		h.log.Debug("(DEV) JSON :: %s", message)
	}
}

/// Atualizando informações sobre a bateria
func (h *WhatsrhymenHandlers) HandleBatteryMessage(message whatsrhymen.BatteryMessage) {
	/* DISABLED FOR NOW
	h.Server.Battery.Timestamp = time.Now()
	h.Server.Battery.Plugged = message.Plugged
	h.Server.Battery.Percentage = message.Percentage
	h.Server.Battery.Powersave = message.Powersave
	*/
}

func (h *WhatsrhymenHandlers) HandleNewContact(contact whatsrhymen.Contact) {
	h.log.Debug("(DEV) NEWCONTACT :: %#v", contact)
}

func (h *WhatsrhymenHandlers) HandleInfoMessage(msg whatsrhymen.MessageInfo) {
	b, err := json.Marshal(msg)
	if err != nil {
		h.log.Error(err)
		return
	}

	h.log.Debug("(DEV) INFO :: %#v", string(b))
}

func (h *WhatsrhymenHandlers) HandleImageMessage(msg whatsrhymen.ImageMessage) {
	message := CreateMessageFromInfo(msg.Info)
	message.Type = whatsapp.ImageMessageType

	FillHeader(message, msg.Info, h.Store())

	//  --> Personalizado para esta seção
	message.Text = "Imagem recebida: " + msg.Type
	FillImageAttachment(message, msg)
	//  <--

	h.Message(message)
}

func (h *WhatsrhymenHandlers) HandleLocationMessage(msg whatsrhymen.LocationMessage) {
	message := CreateMessageFromInfo(msg.Info)
	FillHeader(message, msg.Info, h.Store())

	//  --> Personalizado para esta seção
	message.Text = "Localização recebida ... "
	//  <--

	h.Message(message)
}

func (h *WhatsrhymenHandlers) HandleLiveLocationMessage(msg whatsrhymen.LiveLocationMessage) {
	message := CreateMessageFromInfo(msg.Info)
	FillHeader(message, msg.Info, h.Store())

	//  --> Personalizado para esta seção
	message.Text = "Localização em tempo real recebida ... "
	//  <--

	h.Message(message)
}

func (h *WhatsrhymenHandlers) HandleDocumentMessage(msg whatsrhymen.DocumentMessage) {
	message := CreateMessageFromInfo(msg.Info)
	message.Type = whatsapp.DocumentMessageType

	FillHeader(message, msg.Info, h.Store())

	//  --> Personalizado para esta seção
	innerMSG := msg.Info.Source.Message.DocumentMessage
	message.Text = "Documento recebido: " + msg.Type + " :: " + *innerMSG.Mimetype + " :: " + msg.FileName

	FillDocumentAttachment(message, msg)
	//  <--

	h.Message(message)
}

func (h *WhatsrhymenHandlers) HandleContactMessage(msg whatsrhymen.ContactMessage) {
	message := CreateMessageFromInfo(msg.Info)
	FillHeader(message, msg.Info, h.Store())

	//  --> Personalizado para esta seção
	message.Text = "Contato VCARD recebido ... "
	//  <--

	h.Message(message)
}

func (h *WhatsrhymenHandlers) HandleAudioMessage(msg whatsrhymen.AudioMessage) {
	message := CreateMessageFromInfo(msg.Info)
	message.Type = whatsapp.AudioMessageType
	FillHeader(message, msg.Info, h.Store())

	//  --> Personalizado para esta seção
	message.Text = "Audio recebido: " + msg.Type
	FillAudioAttachment(message, msg)
	//  <--

	h.Message(message)
}

func (h *WhatsrhymenHandlers) HandleTextMessage(msg whatsrhymen.TextMessage) {
	message := CreateMessageFromInfo(msg.Info)
	message.Type = whatsapp.TextMessageType
	FillHeader(message, msg.Info, h.Store())

	//  --> Personalizado para esta seção
	message.Text = msg.Text
	//  <--

	h.Message(message)
}

// Se chamado em modo síncrono, mantém os handlers ativos
func (h *WhatsrhymenHandlers) ShouldCallSynchronously() bool {
	return false
	//return h.Synchronous
}
