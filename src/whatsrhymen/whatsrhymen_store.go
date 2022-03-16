package whatsrhymen

import (
	"database/sql"
	"fmt"

	whatsrhymen "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

const SessionsTable string = "whatsrhymen_sessions"

type WhatsrhymenStore struct {
	WID  string `db:"our_jid"`
	Data []byte `db:"session"`
}

type IWhatsrhymenStore interface {
	Create(wid string) (WhatsrhymenStore, error)
	Get(wid string) (WhatsrhymenStore, error)
	GetOrCreate(wid string) (WhatsrhymenStore, error)
	Update(wid string, data []byte) ([]byte, error)
	Delete(wid string) error
	Exists(wid string) (bool, error)
}

type WhatsrhymenStoreSql struct {
	db     *sql.DB
	logger *log.Logger
}

func NewStore(dialect, address string, logger *log.Logger) (*WhatsrhymenStoreSql, error) {
	db, err := sql.Open(dialect, address)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	container := &WhatsrhymenStoreSql{db, logger}
	err = container.Upgrade()
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade database: %w", err)
	}
	return container, nil
}

func (source WhatsrhymenStoreSql) Exists(wid string) (bool, error) {
	counter := -1
	sqlStatement := "SELECT count(*) FROM whatsrhymen_sessions WHERE Wid = ?"
	row := source.db.QueryRow(sqlStatement, wid)
	err := row.Scan(&counter)
	return counter > 0, err
}

func (source WhatsrhymenStoreSql) Get(wid string) (session whatsrhymen.Session, err error) {
	sqlStatement := "SELECT * FROM whatsrhymen_sessions WHERE Wid = ?;"
	row := source.db.QueryRow(sqlStatement, wid)
	err = row.Scan(&session.Wid, &session.ClientId, &session.ClientToken, &session.ServerToken, &session.EncKey, &session.MacKey)
	return
}

// Create or Insert
func (source WhatsrhymenStoreSql) Update(session whatsrhymen.Session) error {
	exists, err := source.Exists(session.Wid)
	if err != nil {
		return err
	}

	if exists {
		query := "UPDATE whatsrhymen_sessions SET ClientId = ?, ClientToken = ?, ServerToken = ?, EncKey = ?, MacKey = ? WHERE Wid = ?"
		_, err = source.db.Exec(query, session.ClientId, session.ClientToken, session.ServerToken, session.EncKey, session.MacKey, session.Wid)
	} else {
		query := `INSERT INTO whatsrhymen_sessions (Wid, ClientId, ClientToken, ServerToken, EncKey, MacKey) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = source.db.Exec(query, session.Wid, session.ClientId, session.ClientToken, session.ServerToken, session.EncKey, session.MacKey)
	}

	return err
}

func (source WhatsrhymenStoreSql) Delete(wid string) error {
	query := "DELETE FROM whatsrhymen_sessions WHERE Wid = ?"
	_, err := source.db.Exec(query, wid)
	return err
}
