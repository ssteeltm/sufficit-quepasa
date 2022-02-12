package whatsmeow

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
	. "go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type WhatsmeowHandlers struct {
	Client         *Client
	WAHandlers     IWhatsappHandlers
	eventHandlerID uint32
}

func (handler *WhatsmeowHandlers) Register() (err error) {
	if handler.Client.Store == nil {
		err = fmt.Errorf("this client lost the store, probably a logout from whatsapp phone")
		return
	}

	handler.eventHandlerID = handler.Client.AddEventHandler(handler.EventsHandler)
	return
}

// Define os diferentes tipos de eventos a serem reconhecidos
// Aqui se define se vamos processar mensagens | confirmações de leitura | etc
func (handler *WhatsmeowHandlers) EventsHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		handler.Message(v)
		//case *events.Receipt: fmt.Println("Received a receipt! %s", v)

	case *events.LoggedOut:
		log.Error("loggedout ....")
	}
}

//region EVENT MESSAGE

// Aqui se processar um evento de recebimento de uma mensagem genérica
func (handler *WhatsmeowHandlers) Message(evt *events.Message) {

	if evt.Message == nil {
		log.Error("nil message on receiving whatsmeow events | try use rawMessage !")
		return
	}

	b, err := json.Marshal(evt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	message := &WhatsappMessage{Content: evt.Message}

	// basic information
	message.ID = evt.Info.ID
	message.Timestamp = evt.Info.Timestamp
	message.FromMe = evt.Info.IsFromMe

	message.Chat = WhatsappChat{}
	message.Chat.ID = evt.Info.Chat.User

	if evt.Info.IsGroup {
		gInfo, _ := handler.Client.GetGroupInfo(evt.Info.Chat)
		if gInfo != nil {
			message.Chat.Title = gInfo.Name
		}

		message.Participant = WhatsappEndpoint{}
		message.Participant.ID = evt.Info.Sender.User
		message.Participant.Title = evt.Info.PushName
	} else {
		message.Chat.Title = evt.Info.PushName
	}

	message.Text = evt.Message.GetConversation()

	// Process diferent message types
	HandleMessage(message, evt.Message)

	go handler.WAHandlers.Message(message)
}

//endregion
