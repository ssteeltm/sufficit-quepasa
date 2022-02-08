package models

import (
	"log"
	"sync"
)

// Serviço que controla os servidores / bots individuais do whatsapp
type QPWhatsappService struct {
	Servers map[string]*QPWhatsappServer
	DB      *QPDatabase
	Sync    *sync.Mutex // Objeto de sinaleiro para evitar chamadas simultâneas a este objeto
}

var WhatsAppService *QPWhatsappService

func QPWhatsAppStart() {
	log.Println("Starting WhatsApp Service ....")

	servers := make(map[string]*QPWhatsappServer)
	sync := &sync.Mutex{}
	db := *GetDatabase()
	WhatsAppService = &QPWhatsappService{servers, &db, sync}

	// iniciando servidores e cada bot individualmente
	err := WhatsAppService.initService()
	if err != nil {
		log.Printf("Problema ao instanciar bots .... %s", err)
	}
}

// Inclui um novo servidor em um serviço já em andamento
// *Usado quando se passa pela verificação do QRCode
// *Usado quando se inicializa o sistema
func (service *QPWhatsappService) AppendNewServer(bot *QPBot) {
	// Trava simultaneos
	service.Sync.Lock()

	// Cria um novo servidor
	server := NewQPWhatsappServer(bot)

	// Adiciona na lista de servidores
	service.Servers[bot.ID] = server

	// Destrava simultaneos
	service.Sync.Unlock()

	// Inicializa o servidor
	go server.Initialize()
}

// Função privada que irá iniciar todos os servidores apartir do banco de dados
func (service *QPWhatsappService) initService() error {
	bots, err := service.DB.Bot.FindAll()
	if err != nil {
		return err
	}

	for _, bot := range bots {

		if !bot.Verified {

		}
		service.AppendNewServer(bot)
	}

	return nil
}
