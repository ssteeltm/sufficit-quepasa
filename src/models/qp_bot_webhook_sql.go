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
	var result []QpWebhook
	err = source.db.Select(&result, "SELECT url, forwardinternal FROM webhooks WHERE context = ? AND url = ?", context, url)
	if err != nil {
		return
	}

	if len(result) > 0 {
		response = &QpBotWebhook{
			Context:   context,
			QpWebhook: &result[0],
		}
	}
	return
}

func (source QpBotWebhookSql) FindAll(context string) ([]*QpBotWebhook, error) {
	result := []*QpBotWebhook{}
	err := source.db.Select(&result, "SELECT url, forwardinternal FROM webhooks WHERE context = ?", context)
	return result, err
}

func (source QpBotWebhookSql) All() ([]QpBotWebhook, error) {
	result := []QpBotWebhook{}
	err := source.db.Select(&result, "SELECT * FROM webhooks")
	return result, err
}

func (source QpBotWebhookSql) Add(context string, url string, forwardinternal bool) error {
	query := `INSERT OR IGNORE INTO webhooks (context, url, forwardinternal) VALUES (?, ?, ?)`
	_, err := source.db.Exec(query, context, url, forwardinternal)
	return err
}

func (source QpBotWebhookSql) Update(element QpBotWebhook) error {
	query := `UPDATE webhooks SET forwardinternal = ? WHERE context = ? AND url = ?`
	_, err := source.db.Exec(query, element.ForwardInternal, element.Context, element.Url)
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
