package whatsrhymen

import (
	"strings"
	"sync"

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

		WhatsrhymenService = &WhatsrhymenServiceModel{
			Container: container,
		}
	}
}

// Used for scan QR Codes
// Dont forget to attach handlers after success login
func (service *WhatsrhymenServiceModel) CreateEmptyConnection() (conn *WhatsrhymenConnection, err error) {
	logger := log.StandardLogger()
	logger.SetLevel(log.DebugLevel)
	loggerEntry := log.NewEntry(logger)

	conn = &WhatsrhymenConnection{
		Reconnect:      true,
		log:            loggerEntry,
		syncConnection: &sync.Mutex{},
	}

	go conn.EnsureUnderlying()
	return
}

func (service *WhatsrhymenServiceModel) CreateConnection(wid string, logger *log.Logger) (conn *WhatsrhymenConnection, err error) {
	if logger == nil {
		logger = log.StandardLogger()
	}

	logger.SetLevel(log.DebugLevel)
	var loggerEntry *log.Entry
	if len(wid) > 0 {
		loggerEntry = logger.WithField("wid", wid)
	} else {
		loggerEntry = log.NewEntry(logger)
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
		Session:        &session,
		log:            loggerEntry,
		syncConnection: &sync.Mutex{},
		failedToken:    false,
	}

	if len(wid) > 0 {
		go conn.UpdateClient()
	} else {
		go conn.EnsureUnderlying()
	}
	return
}

// Flush entire Whatsrhymen Database
// Use with wisdom !
func (service *WhatsrhymenServiceModel) FlushDatabase() (err error) {
	service.Container.logger.Warn("flushing entire database of whatsrhymen")
	return
}

func (service *WhatsrhymenServiceModel) Delete(wid string) error {
	service.Container.logger.Info("deleting whatsrhymen")
	return service.Container.Delete(wid)
}

func (service *WhatsrhymenServiceModel) UpdateSession(session whatsrhymen.Session) error {
	return service.Container.Update(session)
}
