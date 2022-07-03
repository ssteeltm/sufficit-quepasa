package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type QpBotWebhookSql struct {
	db *sqlx.DB
}

func (source QpBotWebhookSql) Find(context string, url string) (*QpBotWebhook, error) {
	var result *QpBotWebhook
	err := source.db.Get(&result, "SELECT url, failure FROM webhooks WHERE context = ? AND url = ?", context, url)
	return result, err
}

func (source QpBotWebhookSql) FindAll(context string) ([]*QpBotWebhook, error) {
	result := []*QpBotWebhook{}
	err := source.db.Select(&result, "SELECT url, failure FROM webhooks WHERE context = ?", context)
	return result, err
}

func (source QpBotWebhookSql) All() ([]QpBotWebhook, error) {
	result := []QpBotWebhook{}
	err := source.db.Select(&result, "SELECT * FROM webhooks")
	return result, err
}

func (source QpBotWebhookSql) Add(context string, url string) error {
	query := `INSERT OR IGNORE INTO webhooks (context, url, failure) VALUES (?, ?, NULL)`
	_, err := source.db.Exec(query, context, url)
	return err
}

func (source QpBotWebhookSql) Remove(context string, url string) error {
	query := `DELETE FROM webhooks WHERE context = ? AND url = ?`
	_, err := source.db.Exec(query, context, url)
	return err
}

func (source QpBotWebhookSql) Clear(context string) error {
	query := `DELETE FROM webhooks WHERE context = ?`
	_, err := source.db.Exec(query, context)
	return err
}
