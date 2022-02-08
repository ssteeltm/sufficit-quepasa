package whatsmeow

import (
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsmeowServiceModel struct {
	Container *sqlstore.Container
}

var WhatsmeowService *WhatsmeowServiceModel

func WhatsmeowStart() {
	log.Trace("Starting Whatsmeow Service ....")

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:whatsmeow.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	WhatsmeowService = &WhatsmeowServiceModel{Container: container}
}
