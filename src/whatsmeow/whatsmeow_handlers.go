package whatsmeow

import (	
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types/events"
    "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

type WhatsmeowHandlers struct {
    Client              *whatsmeow.Client    
	WAHandlers          whatsapp.IWhatsappHandlers
    eventHandlerID      uint32
}

func (handler *WhatsmeowHandlers) Register() {
    handler.eventHandlerID = handler.Client.AddEventHandler(handler.EventsHandler)
}

// Define os diferentes tipos de eventos a serem reconhecidos
// Aqui se define se vamos processar mensagens | confirmações de leitura | etc
func (handler *WhatsmeowHandlers) EventsHandler(evt interface{}) {
	switch v := evt.(type) {
        case *events.Message: handler.Message(v)
        //case *events.Receipt: fmt.Println("Received a receipt! %s", v)
    }
}

//region EVENT MESSAGE

// Aqui se processar um evento de recebimento de uma mensagem genérica
func (handler *WhatsmeowHandlers) Message(evt *events.Message){
    message := &whatsapp.WhatsappMessage{ Content: evt.Message }

    // basic information
    message.ID = evt.Info.ID
    message.Timestamp = evt.Info.Timestamp
    message.FromMe = evt.Info.IsFromMe
    
    message.Chat = whatsapp.WhatsappChat{}
    message.Chat.ID = evt.Info.Chat.User
    message.Chat.Title = evt.Info.PushName
         
    if evt.Info.IsGroup { 
        message.Participant = whatsapp.WhatsappEndpoint{}
        message.Participant.ID = evt.Info.Sender.User
    }

    message.Text = evt.Message.GetConversation()  

    // Process diferent message types
    HandleMessage(message, evt.Message)

    go handler.WAHandlers.Message(message)
}

//endregion