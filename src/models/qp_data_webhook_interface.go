package models

type QpDataWebhookInterface interface {
	Find(context string, url string) (*QpBotWebhook, error)
	FindAll(context string) ([]*QpBotWebhook, error)
	All() ([]*QpBotWebhook, error)
	Add(element QpBotWebhook) error
	Update(element QpBotWebhook) error
	Remove(context string, url string) error
	Clear(context string) error
}
