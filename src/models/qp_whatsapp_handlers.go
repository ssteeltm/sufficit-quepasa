package models

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// Serviço que controla os servidores / bots individuais do whatsapp
type QPWhatsappHandlers struct {
	messages map[string]WhatsappMessage
	sync     *sync.Mutex // Objeto de sinaleiro para evitar chamadas simultâneas a este objeto

	// Appended events handler
	aeh []interface{ Handle(WhatsappMessage) }

	//filters
	HandleGroups    bool
	HandleBroadcast bool
}

//region CONTRUCTORS

// Create a new QuePasa WhatsApp Event Handler
func NewQPWhatsappHandlers(groups bool, broadcast bool) (handler *QPWhatsappHandlers) {
	handlerMessages := make(map[string]WhatsappMessage)
	handlerSync := &sync.Mutex{}
	handler = &QPWhatsappHandlers{
		messages:        handlerMessages,
		sync:            handlerSync,
		HandleGroups:    groups,
		HandleBroadcast: broadcast,
	}
	return
}

//endregion
//region EVENTS FROM WHATSAPP SERVICE

func (handler *QPWhatsappHandlers) Message(msg *WhatsappMessage) {

	// skipping groups if choosed
	if !handler.HandleGroups && msg.FromGroup() {
		return
	}

	// skipping broadcast if choosed
	if !handler.HandleBroadcast && msg.FromBroadcast() {
		return
	}

	log.Trace("msg recebida/(enviada por outro meio) em models: %s", msg.ID)

	err := handler.appendMsgToCache(msg)
	if err != nil {
		log.Error("error on append msg to cache: %s", err)
	}

}

//endregion

// Salva em cache e inicia gatilhos assíncronos
func (handler *QPWhatsappHandlers) appendMsgToCache(msg *WhatsappMessage) error {

	handler.sync.Lock() // Sinal vermelho para atividades simultâneas
	// Apartir deste ponto só se executa um por vez

	handler.messages[msg.ID] = *msg

	handler.sync.Unlock() // Sinal verde !

	// Executando WebHook de forma assincrona
	handler.Trigger(*msg)

	return nil
}

func (handler *QPWhatsappHandlers) GetMessages(timestamp time.Time) (messages []WhatsappMessage) {
	handler.sync.Lock() // Sinal vermelho para atividades simultâneas
	// Apartir deste ponto só se executa um por vez

	for _, item := range handler.messages {
		if item.Timestamp.After(timestamp) {
			messages = append(messages, item)
		}
	}

	handler.sync.Unlock() // Sinal verde !
	return
}

// Get a single message if exists
func (handler *QPWhatsappHandlers) GetMessage(id string) (msg WhatsappMessage, err error) {
	handler.sync.Lock() // Sinal vermelho para atividades simultâneas
	// Apartir deste ponto só se executa um por vez

	msg, ok := handler.messages[id]
	if !ok {
		err = fmt.Errorf("message not present on handlers (cache)")
	}

	handler.sync.Unlock() // Sinal verde !
	return msg, err
}

//region EVENT HANDLER TO INTERNAL USE, GENERALY TO WEBHOOK

func (handler *QPWhatsappHandlers) Trigger(payload WhatsappMessage) {
	for _, handler := range handler.aeh {
		go handler.Handle(payload)
	}
}

// Register an event handler that triggers on a new message received on cache
func (handler *QPWhatsappHandlers) Register(evt interface{ Handle(WhatsappMessage) }) {
	if !handler.IsRegistered(evt) {
		handler.aeh = append(handler.aeh, evt)
	}
}

func (handler *QPWhatsappHandlers) IsRegistered(evt interface{}) bool {
	for _, v := range handler.aeh {
		if v == evt {
			return true
		}
	}

	return false
}

//endregion
