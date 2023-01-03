package models

import (
	"encoding/json"
	"fmt"
	"time"
)

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

func (source *QpBotWebhook) All() ([]*QpBotWebhook, error) {
	return source.db.All()
}

// passing extra info as json valid or default string
func (source *QpBotWebhook) GetExtraText() string {
	extraJson, err := json.Marshal(&source.Extra)
	if err != nil {
		return fmt.Sprintf("%v", source.Extra)
	} else {
		return string(extraJson)
	}
}

// trying to get interface from json string or default string
func (source *QpBotWebhook) ParseExtra() {
	extraText := fmt.Sprintf("%v", source.Extra)

	var extraJson interface{}
	err := json.Unmarshal([]byte(extraText), &extraJson)
	if err != nil {
		source.Extra = extraText
	} else {
		source.Extra = extraJson
	}
}

func (source *QpBotWebhook) Add(element QpBotWebhook) error {
	return source.db.Add(element)
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
