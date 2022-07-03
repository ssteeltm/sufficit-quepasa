package models

import "time"

// Webhook model
type QpBotWebhook struct {
	Context string `db:"context" json:"context"`

	db QpDataWebhookInterface
	*QpWebhook
}

func (source *QpBotWebhook) Find(context string, url string) (*QpBotWebhook, error) {
	return source.db.Find(context, url)
}

func (source *QpBotWebhook) FindAll(context string) ([]*QpBotWebhook, error) {
	return source.db.FindAll(context)
}

func (source *QpBotWebhook) All() ([]QpBotWebhook, error) {
	return source.db.All()
}

func (source *QpBotWebhook) Add(context string, url string) error {
	return source.db.Add(context, url)
}

func (source *QpBotWebhook) Remove(context string, url string) error {
	return source.db.Remove(context, url)
}

func (source *QpBotWebhook) Clear(context string) error {
	return source.db.Clear(context)
}

// Implement QpWebhookInterface

func (source *QpBotWebhook) GetUrl() string {
	return source.Url
}

func (source *QpBotWebhook) GetFailure() *time.Time {
	return source.Failure
}
