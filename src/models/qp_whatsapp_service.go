package models

import (
	"sync"

	log "github.com/sirupsen/logrus"
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
func (service *QPWhatsappService) AppendNewServer(bot *QPBot) (server *QPWhatsappServer, err error) {
	wid := bot.ID

	// Important for use in update bot info to base
	// Attaching sql store
	if bot.db == nil {
		bot.db = service.DB.Bot
	}

	// Creating a new instance
	server, err = NewQPWhatsappServer(bot, &service.DB.Webhook)
	if err != nil {
		log.Errorf("error on append new server: %s, :: %s", wid, err.Error())
		return
	}

	// Adiciona na lista de servidores
	log.Infof("updating server on cache: %s", wid)
	service.Servers[wid] = server

	// Inicializa o servidor
	if server != nil {
		go server.Initialize()
	}

	return
}

func (service *QPWhatsappService) GetOrCreateServer(currentUserID string, wid string) (server *QPWhatsappServer, err error) {
	log.Debugf("locating server: %s", wid)
	server, ok := service.Servers[wid]
	if !ok {
		log.Debugf("server: %s, not in cache, looking up database", wid)
		bot, err := service.DB.Bot.GetOrCreate(wid, currentUserID)
		if err != nil {
			return nil, err
		}

		// Vinculando base de dados
		if bot.db == nil {
			bot.db = service.DB.Bot
		}

		log.Debugf("server: %s, found", wid)
		server, err = service.AppendNewServer(&bot)
	} else {
		if server.connection != nil && !server.connection.IsInterfaceNil() {
			server.connection.Dispose()
		}
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

	if !service.Initialized {

		bots, err := service.DB.Bot.FindAll()
		if err != nil {
			return err
		}

		for _, bot := range bots {

			// appending server to cache
			service.AppendNewServer(bot)
		}

		service.Initialized = true
	}

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

func (service *QPWhatsappService) GetUser(email string, password string) (user QPUser, err error) {
	log.Debug("finding user ...")
	return service.DB.User.Check(email, password)
}
