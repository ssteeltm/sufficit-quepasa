package whatsrhymen

import (
	"database/sql"
)

type upgradeFunc func(*sql.Tx, *WhatsrhymenStoreSql) error

// Upgrades is a list of functions that will upgrade a database to the latest version.
//
// This may be of use if you want to manage the database fully manually, but in most cases you
// should just call Container.Upgrade to let the library handle everything.
var Upgrades = [...]upgradeFunc{upgradeV1}

func (c *WhatsrhymenStoreSql) getVersion() (int, error) {
	_, err := c.db.Exec("CREATE TABLE IF NOT EXISTS whatsrhymen_version (version INTEGER)")
	if err != nil {
		return -1, err
	}

	version := 0
	row := c.db.QueryRow("SELECT version FROM whatsrhymen_version LIMIT 1")
	if row != nil {
		_ = row.Scan(&version)
	}
	return version, nil
}

func (c *WhatsrhymenStoreSql) setVersion(tx *sql.Tx, version int) error {
	_, err := tx.Exec("DELETE FROM whatsrhymen_version")
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO whatsrhymen_version (version) VALUES ($1)", version)
	return err
}

// Upgrade upgrades the database from the current to the latest version available.
func (c *WhatsrhymenStoreSql) Upgrade() error {
	version, err := c.getVersion()
	if err != nil {
		return err
	}

	for ; version < len(Upgrades); version++ {
		var tx *sql.Tx
		tx, err = c.db.Begin()
		if err != nil {
			return err
		}

		migrateFunc := Upgrades[version]
		err = migrateFunc(tx, c)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if err = c.setVersion(tx, version+1); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func upgradeV1(tx *sql.Tx, _ *WhatsrhymenStoreSql) error {
	_, err := tx.Exec(`CREATE TABLE whatsrhymen_sessions (
		Wid  TEXT,
		ClientId  TEXT,
		ClientToken  TEXT,
		ServerToken  TEXT,
		EncKey  bytea,
		MacKey  bytea,

		PRIMARY KEY (Wid)
	)`)
	if err != nil {
		return err
	}
	return nil
}
