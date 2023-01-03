package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type QpBotWebhookSql struct {
	db *sqlx.DB
}

func (source QpBotWebhookSql) Find(context string, url string) (response *QpBotWebhook, err error) {
	var result []QpBotWebhook
	err = source.db.Select(&result, "SELECT * FROM webhooks WHERE context = ? AND url = ?", context, url)
	if err != nil {
		return
	}

	// adjust extra information
	for _, element := range result {
		element.ParseExtra()
	}

	return
}

func (source QpBotWebhookSql) FindAll(context string) ([]*QpBotWebhook, error) {
	result := []*QpBotWebhook{}
	err := source.db.Select(&result, "SELECT * FROM webhooks WHERE context = ?", context)

	// adjust extra information
	for _, element := range result {
		element.ParseExtra()
	}
	return result, err
}

func (source QpBotWebhookSql) All() ([]*QpBotWebhook, error) {
	result := []*QpBotWebhook{}
	err := source.db.Select(&result, "SELECT * FROM webhooks")

	// adjust extra information
	for _, element := range result {
		element.ParseExtra()
	}
	return result, err
}

func (source QpBotWebhookSql) Add(element QpBotWebhook) error {
	query := `INSERT OR IGNORE INTO webhooks (context, url, forwardinternal, trackid, extra) VALUES (?, ?, ?, ?, ?)`
	_, err := source.db.Exec(query, element.Context, element.Url, element.ForwardInternal, element.TrackId, element.GetExtraText())
	return err
}

func (source QpBotWebhookSql) Update(element QpBotWebhook) error {
	query := `UPDATE webhooks SET forwardinternal = ?, trackid = ?, extra = ? WHERE context = ? AND url = ?`
	_, err := source.db.Exec(query, element.ForwardInternal, element.TrackId, element.GetExtraText(), element.Context, element.Url)
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
