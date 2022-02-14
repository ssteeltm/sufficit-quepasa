package models

import (
	"sync"

	log "github.com/sirupsen/logrus"
	. "github.com/sufficit/sufficit-quepasa-fork/whatsapp"
)

// Serviço que controla os servidores / bots individuais do whatsapp
type QPWhatsappService struct {
	Servers     map[string]*QPWhatsappServer
	DB          *QPDatabase
	Initialized bool

	initlock   *sync.Mutex
	appendlock *sync.Mutex
}

var WhatsappService *QPWhatsappService

func QPWhatsappStart() (err error) {
	if WhatsappService == nil {
		log.Trace("starting whatsapp service ....")

		servers := make(map[string]*QPWhatsappServer)
		db := GetDatabase()
		WhatsappService = &QPWhatsappService{
			Servers:    servers,
			DB:         db,
			initlock:   &sync.Mutex{},
			appendlock: &sync.Mutex{},
		}

		// iniciando servidores e cada bot individualmente
		err = WhatsappService.Initialize()
	} else {
		log.Debug("attempt to start whatsapp service, already started ...")
	}
	return
}

// Inclui um novo servidor em um serviço já em andamento
// *Usado quando se passa pela verificação do QRCode
// *Usado quando se inicializa o sistema
func (service *QPWhatsappService) AppendNewServer(bot *QPBot, con IWhatsappConnection) (server *QPWhatsappServer, err error) {

	// Trava simultaneos
	service.appendlock.Lock()

	// Vinculando base de dados
	bot.db = service.DB.Bot

	// Cria um novo servidor
	server, err = NewQPWhatsappServer(bot, con)
	if err != nil {
		log.Error(err, "error on append new server")
		bot.UpdateVerified(false)
	} else {
		// Adiciona na lista de servidores
		service.Servers[bot.ID] = server
	}

	// Trava simultaneos
	service.appendlock.Unlock()

	// Inicializa o servidor
	if server != nil {
		go server.Initialize()
	}

	return
}

func (service *QPWhatsappService) GetOrCreate(con IWhatsappConnection, userid string) (server *QPWhatsappServer, err error) {
	wid, err := con.GetWid()
	if err != nil {
		return
	}

	server, ok := service.Servers[wid]
	if !ok {
		bot, err := service.DB.Bot.GetOrCreate(wid, userid)
		if err != nil {
			return nil, err
		}
		server, err = service.AppendNewServer(&bot, con)
	}
	return
}

func (service *QPWhatsappService) Delete(server *QPWhatsappServer) (err error) {
	wid := server.GetWid()

	err = server.Delete()
	if err != nil {
		return
	}

	delete(service.Servers, wid)
	return
}

// Função que irá iniciar todos os servidores apartir do banco de dados
func (service *QPWhatsappService) Initialize() (err error) {
	// Trava simultaneos
	service.initlock.Lock()

	if !service.Initialized {

		bots, err := service.DB.Bot.FindAll()
		if err != nil {
			return err
		}

		for _, bot := range bots {
			service.AppendNewServer(bot, nil)
		}

		service.Initialized = true
	}

	// Destrava simultaneos
	service.initlock.Unlock()

	return
}

// Função privada que irá iniciar todos os servidores apartir do banco de dados
func (service *QPWhatsappService) GetServersForUser(userid string) (servers map[string]*QPWhatsappServer) {
	servers = make(map[string]*QPWhatsappServer)
	for _, server := range service.Servers {
		if server.GetOwnerID() == userid {
			servers[server.GetWid()] = server
		}
	}
	return
}
