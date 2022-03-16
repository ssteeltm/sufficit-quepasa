package whatsrhymen

import (
	"strings"
	"time"

	whatsrhymen "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

type WhatsrhymenServiceModel struct {
	Container *WhatsrhymenStoreSql
}

var WhatsrhymenService *WhatsrhymenServiceModel

func (service *WhatsrhymenServiceModel) Start() {
	if service == nil {
		log.Trace("Starting Whatsmeow Service ....")

		dbLog := log.New()
		container, err := NewStore("sqlite3", "file:whatsrhymen.db?_foreign_keys=on", dbLog)
		if err != nil {
			panic(err)
		}

		WhatsrhymenService = &WhatsrhymenServiceModel{Container: container}
	}
}

func (service *WhatsrhymenServiceModel) CreateConnection(wid string, logger *log.Logger) (conn *WhatsrhymenConnection, err error) {
	clientLog := log.New()
	client, err := service.GetWhatsAppClient(wid, clientLog)
	if err != nil {
		return
	}

	logger.SetLevel(log.DebugLevel)
	loggerEntry := logger.WithField("wid", wid)
	handlers := &WhatsrhymenHandlers{
		Connection: client,
		log:        loggerEntry,
	}

	err = handlers.Register()
	if err != nil {
		return
	}

	// Include search for session data here !
	session, err := service.Container.Get(wid)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows in result set") {
			return
		}
		err = nil
	}

	conn = &WhatsrhymenConnection{
		Client:   client,
		Handlers: handlers,
		Session:  &session,
		logger:   logger,
		log:      loggerEntry,
	}
	return
}

/*
func (service *WhatsrhymenServiceModel) GetStoreFromWid(wid string) (str *store.Device, err error) {
	return
}
*/
func (service *WhatsrhymenServiceModel) GetWhatsAppClient(wid string, logger *log.Logger) (client *whatsrhymen.Conn, err error) {
	client, err = whatsrhymen.NewConn(20 * time.Second)

	client.SetClientName("QuePasa for Link", "QuePasa", "0.9")
	client.SetClientVersion(2, 2142, 12)

	log.Printf("debug client version :: %v", client.GetClientVersion())
	return
}

// Flush entire Whatsrhymen Database
// Use with wisdom !
func (service *WhatsrhymenServiceModel) FlushDatabase() (err error) {
	return
}

func (service *WhatsrhymenServiceModel) Delete(wid string) error {
	return service.Container.Delete(wid)
}
