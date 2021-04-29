package models

import (
	"log"
	"sync"
)

// Serviço que controla os servidores / bots individuais do whatsapp
type QPWhatsAppService struct {
	Servers map[string]*QPWhatsAppServer
	Sync    *sync.Mutex // Objeto de sinaleiro para evitar chamadas simultâneas a este objeto
}

var WhatsAppService *QPWhatsAppService

func QPWhatsAppStart() {
	log.Println("Starting WhatsApp Service ....")

	servers := make(map[string]*QPWhatsAppServer)
	sync := &sync.Mutex{}
	WhatsAppService = &QPWhatsAppService{servers, sync}

	// iniciando servidores e cada bot individualmente
	err := WhatsAppService.appendServers()
	if err != nil {
		log.Printf("Problema ao instanciar bots .... %s", err)
	}
}

func (service *QPWhatsAppService) appendServers() error {
	bots, err := FindAllBots(GetDB())
	if err != nil {
		return err
	}

	for _, bot := range bots {
		connection, _ := CreateConnection()

		var handlers *QPMessageHandler

		syncConnetion := &sync.Mutex{}
		syncMessages := &sync.Mutex{}
		recipients := make(map[string]bool)
		messages := make(map[string]QPMessage)
		//var server *QPWhatsAppServer
		server := &QPWhatsAppServer{bot, connection, handlers, recipients, messages, syncConnetion, syncMessages}

		// Adiciona na lista de servidores
		service.Sync.Lock()
		service.Servers[bot.ID] = server
		service.Sync.Unlock()

		go server.Initialize()
	}

	return nil
}