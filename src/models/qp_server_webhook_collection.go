package models

import (
	log "github.com/sirupsen/logrus"
)

// Webhook model
type QpServerWebhookCollection struct {
	Webhooks []*QpWebhook `json:"webhooks,omitempty"`
	context  string
	db       QpDataWebhookInterface
}

// Fill start memory cache
func (source *QpServerWebhookCollection) WebhookFill(context string, db QpDataWebhookInterface) (err error) {
	source.Webhooks = []*QpWebhook{}
	source.context = context
	source.db = db

	whooks, err := source.db.FindAll(source.context)
	if err != nil {
		log.Errorf("erro", err)
		return
	}

	for _, element := range whooks {
		source.Webhooks = append(source.Webhooks, element.QpWebhook)
	}

	return
}

func (source *QpServerWebhookCollection) WebhookAdd(url string) (err error) {
	err = source.db.Add(source.context, url)
	if err == nil {
		newElement := &QpWebhook{
			Url: url,
		}
		exists := false
		for index, element := range source.Webhooks {
			if element.Url == url {
				source.Webhooks = append(source.Webhooks[:index], source.Webhooks[index+1:]...) // remove
				source.Webhooks = append(source.Webhooks, newElement)                           // append a clean one
				exists = true
			}
		}

		if !exists {
			source.Webhooks = append(source.Webhooks, newElement)
		}
	}
	return
}

func (source *QpServerWebhookCollection) WebhookRemove(url string) (err error) {
	err = source.db.Remove(source.context, url)
	if err == nil {
		for index, element := range source.Webhooks {
			if element.Url == url {
				source.Webhooks = append(source.Webhooks[:index], source.Webhooks[index+1:]...)
			}
		}
	}
	return
}

func (source *QpServerWebhookCollection) WebhookClear() (err error) {
	return source.db.Clear(source.context)
}

/*
func (source *QpServerWebhookCollection) WebhookFailure(url string) {
	log.Infof("failure on webhook from: %s", url)
	for index, element := range source.Webhooks {
		if element.Url == url {
			//source.Webhooks[index].Failure = &time.Time{}
		}
	}
}
*/
